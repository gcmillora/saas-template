import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { usePostApiV1Signup } from "@/services/api/v1-public";
import { PasswordRequirements } from "@/components/PasswordRequirements";
import { Loader2 } from "lucide-react";

export function SignUp() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const signup = usePostApiV1Signup();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (password !== confirmPassword) {
      setError("Passwords don't match");
      return;
    }
    try {
      await signup.mutateAsync({
        data: { email, password, confirm_password: confirmPassword },
      });
      navigate("/");
    } catch {
      setError("Signup failed. Please try again.");
    }
  };

  return (
    <div className="flex min-h-screen">
      <div className="flex w-full items-center justify-center px-6 py-12 lg:w-1/2 lg:shrink-0">
        <div className="w-full max-w-sm">
          <div className="mb-8">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">Create your account</h1>
            <p className="mt-1.5 text-sm text-muted-foreground">Get started — it only takes a minute</p>
          </div>
          <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
              <div className="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">{error}</div>
            )}
            <div className="space-y-2">
              <Label htmlFor="email" className="text-sm font-medium">Email</Label>
              <Input id="email" type="email" placeholder="name@example.com" value={email} onChange={(e) => setEmail(e.target.value)} required className="h-11" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password" className="text-sm font-medium">Password</Label>
              <Input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} required minLength={8} className="h-11" />
              <PasswordRequirements password={password} />
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirm-password" className="text-sm font-medium">Confirm Password</Label>
              <Input id="confirm-password" type="password" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} required className="h-11" />
            </div>
            <Button type="submit" className="h-11 w-full text-sm font-medium" disabled={signup.isPending}>
              {signup.isPending ? (<><Loader2 className="mr-2 size-4 animate-spin" />Creating account...</>) : ("Sign Up")}
            </Button>
          </form>
          <p className="mt-8 text-center text-sm text-muted-foreground">
            Already have an account?{" "}
            <Link to="/signin" className="font-medium text-primary transition-colors hover:text-primary/80">Sign in</Link>
          </p>
        </div>
      </div>
      <div className="hidden lg:flex lg:w-1/2 lg:items-center lg:justify-center bg-gradient-to-br from-primary/5 via-primary/10 to-primary/5">
        <div className="h-64 w-64 rounded-full bg-primary/5" />
      </div>
    </div>
  );
}
