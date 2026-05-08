import { useState } from "react";
import { message, Modal, Input, Progress } from "antd";
import { generate3DAndWait } from "../../api/threeD";
import ModelViewer from "../model/modelViewer";
import "./generate3D.scss";

interface IGenerate3DParams {
  onGenerated?: (modelUrl: string, prompt: string) => void;
}

export default function Generate3DButton(params: IGenerate3DParams) {
  const [visible, setVisible] = useState(false);
  const [prompt, setPrompt] = useState("");
  const [loading, setLoading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [statusText, setStatusText] = useState("");
  const [modelUrl, setModelUrl] = useState<string | null>(null);

  const handleGenerate = async () => {
    if (!prompt.trim()) {
      message.warning("请输入3D模型描述");
      return;
    }

    setLoading(true);
    setProgress(0);
    setStatusText("提交中...");
    setModelUrl(null);

    try {
      const files = await generate3DAndWait(prompt, (status, elapsed) => {
        const pct = Math.min(95, Math.round((elapsed / 180) * 100));
        setProgress(pct);
        setStatusText(`生成中... ${elapsed}秒`);
      });

      if (files && files.length > 0) {
        const glbFile = files.find((f) => f.type === "GLB") || files[0];
        setModelUrl(glbFile.url);
        setProgress(100);
        setStatusText("生成完成!");
        params.onGenerated?.(glbFile.url, prompt);
        message.success("3D模型生成完成!");
      }
    } catch (err: any) {
      message.error(err.message || "生成失败");
      setStatusText("生成失败");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <button className="generate3DBtn" onClick={() => setVisible(true)}>
        🎨 3D
      </button>

      <Modal
        title="生成3D模型"
        visible={visible}
        onCancel={() => { setVisible(false); setModelUrl(null); }}
        footer={null}
        width={450}
      >
        <Input.TextArea
          placeholder="描述你想要的3D角色，如：一个卡通风格的战士角色"
          value={prompt}
          onChange={(e) => setPrompt(e.target.value)}
          rows={3}
          disabled={loading}
        />

        <button
          className="generateBtn"
          onClick={handleGenerate}
          disabled={loading}
        >
          {loading ? "生成中..." : "开始生成"}
        </button>

        {loading && (
          <div className="progressBar">
            <Progress percent={progress} status="active" />
            <div className="statusText">{statusText}</div>
          </div>
        )}

        {modelUrl && <ModelViewer modelUrl={modelUrl} width={400} height={350} />}
      </Modal>
    </>
  );
}
