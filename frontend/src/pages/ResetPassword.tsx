import { useState } from "react";
import { Link, useSearchParams, useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PasswordRequirements } from "@/components/PasswordRequirements";
import { usePostApiV1ResetPassword } from "@/services/api/v1-public";
import { Loader2, CheckCircle2, AlertCircle } from "lucide-react";

export function ResetPassword() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token");
  const navigate = useNavigate();

  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const resetPassword = usePostApiV1ResetPassword();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (password !== confirmPassword) {
      setError("Passwords don't match");
      return;
    }
    try {
      await resetPassword.mutateAsync({
        data: { token: token!, password, confirm_password: confirmPassword },
      });
      setSuccess(true);
      setTimeout(() => navigate("/signin"), 3000);
    } catch {
      setError("Invalid or expired reset token. Please request a new one.");
    }
  };

  return (
    <div className="flex min-h-screen">
      <div className="flex w-full items-center justify-center px-6 py-12 lg:w-1/2 lg:shrink-0">
        <div className="w-full max-w-sm">
          {!token ? (
            <>
              <div className="mb-8">
                <div className="mb-4 flex size-12 items-center justify-center rounded-xl bg-destructive/10">
                  <AlertCircle className="size-6 text-destructive" />
                </div>
                <h1 className="text-2xl font-bold tracking-tight text-foreground">Invalid reset link</h1>
                <p className="mt-1.5 text-sm text-muted-foreground">This password reset link is invalid or has expired.</p>
              </div>
              <Button asChild className="h-11 w-full text-sm font-medium">
                <Link to="/forgot-password">Request a new reset link</Link>
              </Button>
            </>
          ) : success ? (
            <>
              <div className="mb-8">
                <div className="mb-4 flex size-12 items-center justify-center rounded-xl bg-emerald-500/10">
                  <CheckCircle2 className="size-6 text-emerald-600" />
                </div>
                <h1 className="text-2xl font-bold tracking-tight text-foreground">Password reset</h1>
                <p className="mt-1.5 text-sm text-muted-foreground">Your password has been reset successfully. Redirecting to sign in...</p>
              </div>
              <Button asChild variant="outline" className="h-11 w-full text-sm font-medium">
                <Link to="/signin">Go to sign in</Link>
              </Button>
            </>
          ) : (
            <>
              <div className="mb-8">
                <h1 className="text-2xl font-bold tracking-tight text-foreground">Reset password</h1>
                <p className="mt-1.5 text-sm text-muted-foreground">Enter your new password below</p>
              </div>
              <form onSubmit={handleSubmit} className="space-y-5">
                {error && (
                  <div className="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">{error}</div>
                )}
                <div className="space-y-2">
                  <Label htmlFor="password" className="text-sm font-medium">New Password</Label>
                  <Input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} required minLength={8} autoFocus className="h-11" />
                  <PasswordRequirements password={password} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="confirm-password" className="text-sm font-medium">Confirm Password</Label>
                  <Input id="confirm-password" type="password" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} required className="h-11" />
                </div>
                <Button type="submit" className="h-11 w-full text-sm font-medium" disabled={resetPassword.isPending}>
                  {resetPassword.isPending ? (<><Loader2 className="mr-2 size-4 animate-spin" />Resetting...</>) : ("Reset Password")}
                </Button>
              </form>
            </>
          )}
        </div>
      </div>
      <div className="hidden lg:flex lg:w-1/2 lg:items-center lg:justify-center bg-gradient-to-br from-primary/5 via-primary/10 to-primary/5">
        <div className="h-64 w-64 rounded-full bg-primary/5" />
      </div>
    </div>
  );
}
