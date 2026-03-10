import { useState } from "react";
import { useNavigate, Navigate } from "react-router-dom";
import { Form, Input, Button, Card, message } from "antd";
import { UserOutlined, LockOutlined } from "@ant-design/icons";
import { useAuth } from "@/lib/auth";
import { getErrorMessage } from "@/lib/api";

export default function LoginPage() {
    const navigate = useNavigate();
    const { isAuthenticated, login } = useAuth();
    const [loading, setLoading] = useState(false);

    if (isAuthenticated) {
        return <Navigate to="/" replace />;
    }

    const handleSubmit = async (values: { identifier: string; password: string }) => {
        setLoading(true);
        try {
            await login(values.identifier, values.password);
            message.success("Đăng nhập thành công");
            navigate("/");
        } catch (error) {
            message.error(getErrorMessage(error));
        } finally {
            setLoading(false);
        }
    };

    return (
        <div
            style={{
                minHeight: "100vh",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)",
            }}
        >
            <Card
                style={{ width: "100%", maxWidth: 400, margin: 16 }}
                title={
                    <div style={{ textAlign: "center", fontSize: 24, fontWeight: 700, color: "#EE4D2D" }}>
                        Golf Store
                    </div>
                }
            >
                <Form name="login" onFinish={handleSubmit} layout="vertical" size="large" data-testid="customer-login-form">
                    <Form.Item
                        name="identifier"
                        rules={[{ required: true, message: "Vui lòng nhập email hoặc số điện thoại" }]}
                    >
                        <Input
                            data-testid="customer-login-identifier"
                            prefix={<UserOutlined />}
                            placeholder="Email hoặc số điện thoại"
                            autoComplete="username"
                        />
                    </Form.Item>

                    <Form.Item name="password" rules={[{ required: true, message: "Vui lòng nhập mật khẩu" }]}>
                        <Input.Password
                            data-testid="customer-login-password"
                            prefix={<LockOutlined />}
                            placeholder="Mật khẩu"
                            autoComplete="current-password"
                        />
                    </Form.Item>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" loading={loading} block data-testid="customer-login-submit">
                            Đăng nhập
                        </Button>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    );
}
