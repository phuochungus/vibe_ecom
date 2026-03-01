import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { Eye, EyeOff, Store, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuth } from "@/lib/auth";
import { getErrorMessage } from "@/lib/api";

interface LoginForm {
    identifier: string;
    password: string;
}

export default function LoginPage() {
    const { login } = useAuth();
    const navigate = useNavigate();
    const [showPassword, setShowPassword] = useState(false);
    const [loading, setLoading] = useState(false);

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<LoginForm>();

    const onSubmit = async (data: LoginForm) => {
        try {
            setLoading(true);
            await login(data.identifier, data.password);
            navigate("/");
        } catch (err) {
            const msg = getErrorMessage(err);
            if (msg.includes("ACCOUNT_LOCKED") || msg.toLowerCase().includes("locked")) {
                toast.error("Tài khoản đã bị khóa tạm 15 phút do đăng nhập sai quá nhiều lần.");
            } else {
                toast.error("Sai email/mật khẩu hoặc tài khoản không hợp lệ.");
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 p-4">
            {/* Background decoration */}
            <div className="absolute inset-0 overflow-hidden pointer-events-none">
                <div className="absolute -top-1/2 -left-1/2 w-full h-full rounded-full bg-primary/5 blur-3xl" />
                <div className="absolute -bottom-1/2 -right-1/2 w-full h-full rounded-full bg-primary/5 blur-3xl" />
            </div>

            <Card className="w-full max-w-md relative z-10 border-white/10 bg-card/80 backdrop-blur shadow-2xl">
                <CardHeader className="space-y-4 pb-6">
                    <div className="flex items-center justify-center">
                        <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-primary text-primary-foreground shadow-lg">
                            <Store className="h-7 w-7" />
                        </div>
                    </div>
                    <div className="text-center space-y-1">
                        <CardTitle className="text-2xl font-bold">Golf Store Admin</CardTitle>
                        <CardDescription>Đăng nhập để quản lý hệ thống</CardDescription>
                    </div>
                </CardHeader>

                <CardContent>
                    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="identifier">Email / Số điện thoại</Label>
                            <Input
                                id="identifier"
                                placeholder="admin@example.com"
                                autoComplete="username"
                                {...register("identifier", { required: "Vui lòng nhập email hoặc số điện thoại" })}
                            />
                            {errors.identifier && (
                                <p className="text-xs text-destructive">{errors.identifier.message}</p>
                            )}
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="password">Mật khẩu</Label>
                            <div className="relative">
                                <Input
                                    id="password"
                                    type={showPassword ? "text" : "password"}
                                    placeholder="••••••••"
                                    autoComplete="current-password"
                                    className="pr-10"
                                    {...register("password", { required: "Vui lòng nhập mật khẩu" })}
                                />
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="sm"
                                    className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                                    onClick={() => setShowPassword((v) => !v)}
                                    tabIndex={-1}
                                >
                                    {showPassword ? (
                                        <EyeOff className="h-4 w-4 text-muted-foreground" />
                                    ) : (
                                        <Eye className="h-4 w-4 text-muted-foreground" />
                                    )}
                                </Button>
                            </div>
                            {errors.password && <p className="text-xs text-destructive">{errors.password.message}</p>}
                        </div>

                        <Button type="submit" className="w-full" disabled={loading}>
                            {loading ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Đang đăng nhập...
                                </>
                            ) : (
                                "Đăng nhập"
                            )}
                        </Button>
                    </form>
                </CardContent>
            </Card>

            {/* Sonner needs to be here for the login page outside the layout */}
            <div id="sonner-login" />
        </div>
    );
}
