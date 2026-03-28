#!/usr/bin/env python3
"""
混元3D API 工具 - 腾讯云混元生3D专业版
用法:
  python3 hunyuan3d.py --prompt "一只小猫" --format GLB
  python3 hunyuan3d.py --query JOB_ID
"""
import hashlib, hmac, json, time, urllib.request, urllib.error, argparse, os, sys
from datetime import datetime, timezone

# 从凭证文件或环境变量读取
CREDS_FILE = os.path.expanduser("~/.openclaw/.tencent_credentials")
creds = {}
if os.path.exists(CREDS_FILE):
    with open(CREDS_FILE) as f:
        for line in f:
            if "=" in line:
                k, v = line.strip().split("=", 1)
                creds[k] = v

SECRET_ID = os.environ.get("TENCENT_SECRET_ID", creds.get("TENCENT_SECRET_ID", ""))
SECRET_KEY = os.environ.get("TENCENT_SECRET_KEY", creds.get("TENCENT_SECRET_KEY", ""))
HOST = "ai3d.tencentcloudapi.com"
SERVICE = "ai3d"
REGION = "ap-guangzhou"
VERSION = "2025-05-13"

def _sign(key, msg):
    return hmac.new(key, msg.encode("utf-8"), hashlib.sha256).digest()

def tc3_request(action, params):
    """发送腾讯云 TC3-HMAC-SHA256 签名请求"""
    timestamp = int(time.time())
    date = datetime.fromtimestamp(timestamp, tz=timezone.utc).strftime("%Y-%m-%d")
    payload = json.dumps(params)
    hashed_payload = hashlib.sha256(payload.encode("utf-8")).hexdigest()
    canonical = f"POST\n/\n\ncontent-type:application/json\nhost:{HOST}\n\ncontent-type;host\n{hashed_payload}"
    scope = f"{date}/{SERVICE}/tc3_request"
    hashed_canonical = hashlib.sha256(canonical.encode("utf-8")).hexdigest()
    to_sign = f"TC3-HMAC-SHA256\n{timestamp}\n{scope}\n{hashed_canonical}"
    
    sd = _sign(("TC3" + SECRET_KEY).encode("utf-8"), date)
    ss = _sign(sd, SERVICE)
    sk = _sign(ss, "tc3_request")
    sig = hmac.new(sk, to_sign.encode("utf-8"), hashlib.sha256).hexdigest()
    
    auth = f"TC3-HMAC-SHA256 Credential={SECRET_ID}/{scope}, SignedHeaders=content-type;host, Signature={sig}"
    headers = {
        "Authorization": auth, "Content-Type": "application/json",
        "Host": HOST, "X-TC-Action": action, "X-TC-Version": VERSION,
        "X-TC-Timestamp": str(timestamp), "X-TC-Region": REGION,
    }
    req = urllib.request.Request(f"https://{HOST}/", data=payload.encode("utf-8"), headers=headers)
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            return json.loads(resp.read())
    except urllib.error.HTTPError as e:
        return json.loads(e.read().decode())

def submit_job(prompt=None, image_url=None, result_format="GLB"):
    """提交混元3D生成任务"""
    params = {"ResultFormat": result_format}
    if prompt:
        params["Prompt"] = prompt
    elif image_url:
        params["ImageUrl"] = image_url
    else:
        raise ValueError("prompt 或 image_url 必填其一")
    
    result = tc3_request("SubmitHunyuanTo3DJob", params)
    resp = result.get("Response", {})
    err = resp.get("Error")
    if err:
        raise Exception(f"{err['Code']}: {err['Message']}")
    return resp["JobId"]

def query_job(job_id):
    """查询任务状态"""
    result = tc3_request("QueryHunyuanTo3DJob", {"JobId": job_id})
    resp = result.get("Response", {})
    err = resp.get("Error")
    if err:
        raise Exception(f"{err['Code']}: {err['Message']}")
    return resp

def generate_and_wait(prompt=None, image_url=None, result_format="GLB", 
                      timeout=300, poll_interval=10):
    """提交任务并等待完成，返回结果"""
    job_id = submit_job(prompt=prompt, image_url=image_url, 
                        result_format=result_format)
    print(f"📋 JobId: {job_id}")
    
    waited = 0
    while waited < timeout:
        time.sleep(poll_interval)
        waited += poll_interval
        resp = query_job(job_id)
        status = resp.get("Status", "UNKNOWN")
        print(f"⏳ [{waited}s] {status}")
        
        if status == "DONE":
            files = resp.get("ResultFile3Ds", [])
            return {"status": "success", "job_id": job_id, "files": files}
        elif status in ("FAIL", "FAILED"):
            return {"status": "failed", "job_id": job_id, 
                    "error": resp.get("ErrorMessage", "unknown")}
    
    return {"status": "timeout", "job_id": job_id}

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="混元3D生成工具")
    parser.add_argument("--prompt", "-p", help="文生3D提示词")
    parser.add_argument("--image", "-i", help="图生3D图片URL")
    parser.add_argument("--format", "-f", default="GLB", choices=["OBJ", "GLB", "STL", "USDZ", "FBX", "MP4"])
    parser.add_argument("--model", "-m", default="3.0", choices=["3.0", "3.1"])
    parser.add_argument("--query", "-q", help="查询已有任务ID")
    parser.add_argument("--timeout", "-t", type=int, default=300)
    args = parser.parse_args()
    
    if args.query:
        resp = query_job(args.query)
        print(json.dumps(resp, indent=2, ensure_ascii=False))
    elif args.prompt or args.image:
        print("🚀 提交混元3D生成任务...")
        result = generate_and_wait(
            prompt=args.prompt, image_url=args.image,
            result_format=args.format, timeout=args.timeout
        )
        print(f"\n{'✅' if result['status']=='success' else '❌'} {result['status']}")
        if result.get("files"):
            for f in result["files"]:
                print(f"  📦 {f.get('Type')}: {f.get('Url', '')[:200]}")
    else:
        parser.print_help()
