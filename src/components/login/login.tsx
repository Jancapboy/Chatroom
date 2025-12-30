import { useState, useEffect } from "react";
import InputBars from "./inputBars";
import "./login.scss";
import { message, Spin, Checkbox } from "antd";
import axios from "axios";
import serverHost from "../../api/hostname";
import IErrRes from "../../types/IErrRes";
import IUserInfo from "../../types/IUserInfo";

interface ILoginParams {
  logged: boolean;
  logAs: (userInfo: IUserInfo) => void;
  setToken: (token: string) => void;
}

// 存储登录信息的键名
const LOGIN_STORAGE_KEY = "chatroom_login_info";

interface IStoredLoginInfo {
  username: string;
  password: string;
  remember: boolean;
}

function Login(params: ILoginParams) {
  const [LogRegister, setLogRegister] = useState<"login" | "register">("login");
  const [waiting, setWaiting] = useState<boolean>(false);
  const [username, setUsername] = useState<string>("");
  const [nickname, setNickname] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [usernameErr, setUsernameErr] = useState<boolean>(false);
  const [nicknameErr, setNicknameErr] = useState<boolean>(false);
  const [passwordErr, setPasswordErr] = useState<boolean>(false);
  const [remember, setRemember] = useState<boolean>(false); // 添加记住密码状态

  // 组件加载时检查本地存储的登录信息
  useEffect(() => {
    const storedInfo = localStorage.getItem(LOGIN_STORAGE_KEY);
    if (storedInfo) {
      try {
        const parsedInfo: IStoredLoginInfo = JSON.parse(storedInfo);
        setUsername(parsedInfo.username);
        setPassword(parsedInfo.password);
        setRemember(parsedInfo.remember);
      } catch (error) {
        console.error("解析本地存储的登录信息失败:", error);
      }
    }
  }, []);

  const checkInputIllegal = (): boolean => {
    let inputIllegal: boolean = false;
    if (!/^(\d|\w|_){5,30}$/g.test(username)) {
      // 检验用户名是否合法
      message.error("用户名不合法:请输入 5-30 位字符", 2);
      inputIllegal = true;
      setUsernameErr(true);
    }
    if (
        // 检验用户昵称是否合法
        LogRegister === "register" &&
        (nickname.length < 1 || nickname.length > 30)
    ) {
      message.error("用户昵称不合法:请输入 1-30 位字符", 2);
      inputIllegal = true;
      setNicknameErr(true);
    }
    if (!/(\d|\w){6,20}$/g.test(password)) {
      // 检验密码输入是否合法
      message.error("密码输入不合法:请输入 6-20 位的数字或大小写字母", 2);
      inputIllegal = true;
      setPasswordErr(true);
    }
    if (waiting) {
      message.warning("正在等待服务器响应", 1);
      inputIllegal = true;
    }

    return inputIllegal;
  };

  const saveLoginInfo = () => {
    if (remember) {
      // 保存登录信息到本地存储
      const loginInfo: IStoredLoginInfo = {
        username,
        password,
        remember: true
      };
      localStorage.setItem(LOGIN_STORAGE_KEY, JSON.stringify(loginInfo));
    } else {
      // 如果取消勾选记住密码，则清除存储的登录信息
      localStorage.removeItem(LOGIN_STORAGE_KEY);
    }
  };

  const launchLogRequest = async (type: "login" | "register") => {
    if (checkInputIllegal()) {
      return;
    }

    // 只有在登录时才保存登录信息
    if (type === "login" && remember) {
      saveLoginInfo();
    }

    let postData =
        LogRegister === "login"
            ? {
              user_name: username,
              password: password,
            }
            : {
              user_name: username,
              nickname: nickname,
              password: password,
            };
    setWaiting(true);
    axios
        .post(serverHost + `/${LogRegister}`, postData)
        .then((res) => {
          if (type === "register") {
            // 注册
            message.success("注册成功，请登录", 2);
            setLogRegister("login");
            // 注册成功后清除本地存储的密码信息
            localStorage.removeItem(LOGIN_STORAGE_KEY);
          } else {
            // 登录
            let data = res.data as {
              user_id: number;
              nickname: string;
              token: string;
            };
            message.success("登陆成功", 2);
            console.log(data);
            params.logAs({
              user_id: data.user_id,
              user_name: username,
              nickname: data.nickname,
            });
            params.setToken(data.token);
          }
        })
        .catch((err) => {
          console.log(err);
          if (!err.response) {
            message.error("服务器错误", 2);
          } else {
            let data = err.response.data as IErrRes;
            message.error(data.msg, 2);
          }
        })
        .finally(() => {
          setWaiting(false);
        });
  };

  return (
      <div className={`login ${params.logged ? "logged" : ""}`}>
        <div className="loginBar">
          <div className="info">
            <div className="title">
              {LogRegister === "login" ? "登 录" : "注 册"}
            </div>
            <InputBars
                LogRegister={LogRegister}
                username={username}
                password={password}
                nickname={nickname}
                setUsername={setUsername}
                setNickname={setNickname}
                setPassword={setPassword}
                usernameErr={usernameErr}
                nicknameErr={nicknameErr}
                passwordErr={passwordErr}
                launchLogRequest={launchLogRequest}
            />

            {/* 只在登录界面显示记住密码选项 */}
            {LogRegister === "login" && (
                <div className="remember-password">
                  <Checkbox
                      checked={remember}
                      onChange={(e) => setRemember(e.target.checked)}
                  >
                    记住密码
                  </Checkbox>
                </div>
            )}

            <div
                className="changeType"
                onClick={() => {
                  if (LogRegister === "login") setLogRegister("register");
                  else setLogRegister("login");
                }}
            >
              {LogRegister === "login"
                  ? "没有账号？点击注册"
                  : "已有账号？点击登陆"}
            </div>
          </div>
          <div
              className="logButton"
              onClick={() => {
                launchLogRequest(LogRegister);
              }}
          >
            {waiting ? (
                <Spin size="large" />
            ) : LogRegister === "login" ? (
                <div className="iconfont">&#xe6dd;</div>
            ) : (
                <div className="iconfont">&#xe683;</div>
            )}
          </div>
        </div>
      </div>
  );
}

export default Login;
