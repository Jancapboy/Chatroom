import hostname from "./hostname";

// 提交3D生成任务
export async function generate3D(
  prompt: string,
  resultFormat: string = "GLB"
): Promise<{ code: number; msg: string; job_id?: string }> {
  const res = await fetch(`${hostname}/api/3d/generate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ prompt, result_format: resultFormat }),
  });
  return res.json();
}

// 查询3D生成结果
export async function query3D(jobId: string): Promise<{
  code: number;
  status: string;
  files?: { type: string; url: string }[];
}> {
  const res = await fetch(`${hostname}/api/3d/query`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ job_id: jobId }),
  });
  return res.json();
}

// 提交并轮询等待结果
export async function generate3DAndWait(
  prompt: string,
  onProgress?: (status: string, elapsed: number) => void,
  timeoutMs: number = 300000
): Promise<{ type: string; url: string }[] | null> {
  const submitRes = await generate3D(prompt);
  if (submitRes.code !== 0 || !submitRes.job_id) {
    throw new Error(submitRes.msg || "提交失败");
  }

  const jobId = submitRes.job_id;
  const startTime = Date.now();

  while (Date.now() - startTime < timeoutMs) {
    await new Promise((r) => setTimeout(r, 10000)); // 10秒轮询
    const elapsed = Math.round((Date.now() - startTime) / 1000);

    const queryRes = await query3D(jobId);
    onProgress?.(queryRes.status, elapsed);

    if (queryRes.status === "DONE" && queryRes.files) {
      return queryRes.files;
    }
    if (queryRes.status === "FAIL" || queryRes.status === "FAILED") {
      throw new Error("生成失败");
    }
  }

  throw new Error("生成超时");
}
