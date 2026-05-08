#!/usr/bin/env python3
"""
Chatroom 3D 集成 - 端到端测试
测试后端API + 前端编译 + 联调
"""
import urllib.request, urllib.error, json, sys, time, subprocess, os

BASE = "http://localhost:4002"
PASS = 0
FAIL = 0

def test(name, func):
    global PASS, FAIL
    try:
        result = func()
        if result:
            print(f"  ✅ {name}")
            PASS += 1
        else:
            print(f"  ❌ {name} - 返回 False")
            FAIL += 1
    except Exception as e:
        print(f"  ❌ {name} - {e}")
        FAIL += 1

def api_post(path, data=None):
    url = BASE + path
    body = json.dumps(data).encode() if data else None
    req = urllib.request.Request(url, data=body, headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req, timeout=15) as resp:
        return json.loads(resp.read())

def api_get(path):
    with urllib.request.urlopen(BASE + path, timeout=10) as resp:
        return json.loads(resp.read())

# ===== 1. 后端 API 测试 =====
print("\n🔧 1. 后端 API 测试")

def test_health():
    r = api_get("/health")
    return r.get("status") == "ok"

def test_generate_missing_prompt():
    """应该返回400"""
    req = urllib.request.Request(BASE + "/api/3d/generate",
        data=b'{}', headers={"Content-Type": "application/json"})
    try:
        urllib.request.urlopen(req, timeout=10)
        return False
    except urllib.error.HTTPError as e:
        return e.code == 400

def test_generate_success():
    r = api_post("/api/3d/generate", {"prompt": "测试用小猫", "result_format": "GLB"})
    return r.get("code") == 0 and r.get("job_id") is not None

def test_query_missing_jobid():
    """应该返回400"""
    req = urllib.request.Request(BASE + "/api/3d/query",
        data=b'{}', headers={"Content-Type": "application/json"})
    try:
        urllib.request.urlopen(req, timeout=10)
        return False
    except urllib.error.HTTPError as e:
        return e.code == 400

def test_query_success():
    """用之前成功的JobId查询"""
    r = api_post("/api/3d/query", {"job_id": "1429320300369436672"})
    return r.get("code") == 0 and r.get("status") in ("DONE", "RUN", "WAIT")

def test_cors_preflight():
    """OPTIONS请求应返回CORS头"""
    req = urllib.request.Request(BASE + "/api/3d/generate", method="OPTIONS")
    req.add_header("Origin", "http://localhost:3000")
    with urllib.request.urlopen(req, timeout=5) as resp:
        return resp.status == 204

test("健康检查", test_health)
test("生成-缺少prompt应400", test_generate_missing_prompt)
test("生成-正常提交", test_generate_success)
test("查询-缺少job_id应400", test_query_missing_jobid)
test("查询-正常查询", test_query_success)
test("CORS预检请求", test_cors_preflight)

# ===== 2. 前端编译测试 =====
print("\n🎨 2. 前端代码检查")

def test_frontend_files_exist():
    files = [
        "/tmp/chatroom-frontend/src/components/model/modelViewer.tsx",
        "/tmp/chatroom-frontend/src/components/model/generate3D.tsx",
        "/tmp/chatroom-frontend/src/api/threeD.ts",
    ]
    return all(os.path.exists(f) for f in files)

def test_three_js_installed():
    return os.path.exists("/tmp/chatroom-frontend/node_modules/three/build/three.module.js")

def test_frontend_build_exists():
    return os.path.exists("/tmp/chatroom-frontend/build/index.html")

test("前端3D文件存在", test_frontend_files_exist)
test("three.js已安装", test_three_js_installed)
test("前端build产物存在", test_frontend_build_exists)

# ===== 3. 后端Go代码检查 =====
print("\n📦 3. 后端代码检查")

def test_backend_files_exist():
    files = [
        "/tmp/chatroom-3d/internal/service/hunyuan3d.go",
        "/tmp/chatroom-3d/internal/routers/api/3d.go",
        "/tmp/chatroom-3d/internal/routers/routers.go",
        "/tmp/chatroom-3d/cmd/test_3d_server/main.go",
    ]
    return all(os.path.exists(f) for f in files)

def test_binary_exists():
    return os.path.exists("/tmp/test_3d_server")

test("后端Go文件存在", test_backend_files_exist)
test("后端二进制存在", test_binary_exists)

# ===== 4. Git状态检查 =====
print("\n📋 4. Git 状态检查")

def test_git_branch():
    r = subprocess.run(["git", "-C", "/tmp/chatroom-3d", "branch", "--show-current"],
                       capture_output=True, text=True)
    return r.stdout.strip() == "feature/3d-integration"

def test_git_clean():
    r = subprocess.run(["git", "-C", "/tmp/chatroom-3d", "status", "--porcelain"],
                       capture_output=True, text=True)
    # 只检查tracked文件是否有未提交修改
    lines = [l for l in r.stdout.strip().split("\n") if l and not l.startswith("??")]
    return len(lines) == 0

test("Git分支正确", test_git_branch)
test("Git无未提交修改", test_git_clean)

# ===== 结果汇总 =====
total = PASS + FAIL
print(f"\n{'='*40}")
print(f"📊 测试结果: {PASS}/{total} 通过")
if FAIL > 0:
    print(f"   ❌ {FAIL} 个失败")
    sys.exit(1)
else:
    print(f"   ✅ 全部通过!")
    sys.exit(0)
