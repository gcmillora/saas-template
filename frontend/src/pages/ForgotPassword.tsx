import { useState } from "react";
import { Link } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { usePostApiV1ForgotPassword } from "@/services/api/v1-public";
import { Loader2, ArrowLeft, MailCheck } from "lucide-react";

export function ForgotPassword() {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const forgotPassword = usePostApiV1ForgotPassword();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await forgotPassword.mutateAsync({ data: { email } });
    } catch {
      // Always show success to prevent email enumeration
    }
    setSubmitted(true);
  };

  return (
    <div className="flex min-h-screen">
      <div className="flex w-full items-center justify-center px-6 py-12 lg:w-1/2 lg:shrink-0">
        <div className="w-full max-w-sm">
          {submitted ? (
            <>
              <div className="mb-8">
                <div className="mb-4 flex size-12 items-center justify-center rounded-xl bg-primary/10">
                  <MailCheck className="size-6 text-primary" />
                </div>
                <h1 className="text-2xl font-bold tracking-tight text-foreground">
                  Check your email
                </h1>
                <p className="mt-1.5 text-sm text-muted-foreground">
                  If an account with that email exists, we've sent a password
                  reset link. Check your inbox.
                </p>
              </div>
              <Button asChild variant="outline" className="h-11 w-full text-sm font-medium">
                <Link to="/signin">
                  <ArrowLeft className="size-4" />
                  Back to sign in
                </Link>
              </Button>
            </>
          ) : (
            <>
              <div className="mb-8">
                <h1 className="text-2xl font-bold tracking-tight text-foreground">
                  Forgot password?
                </h1>
                <p className="mt-1.5 text-sm text-muted-foreground">
                  Enter your email and we'll send you a reset link
                </p>
              </div>
              <form onSubmit={handleSubmit} className="space-y-5">
                <div className="space-y-2">
                  <Label htmlFor="email" className="text-sm font-medium">Email</Label>
                  <Input id="email" type="email" placeholder="name@example.com" value={email} onChange={(e) => setEmail(e.target.value)} required autoFocus className="h-11" />
                </div>
                <Button type="submit" className="h-11 w-full text-sm font-medium" disabled={forgotPassword.isPending}>
                  {forgotPassword.isPending ? (<><Loader2 className="mr-2 size-4 animate-spin" />Sending...</>) : ("Send Reset Link")}
                </Button>
              </form>
              <p className="mt-8 text-center text-sm text-muted-foreground">
                Remember your password?{" "}
                <Link to="/signin" className="font-medium text-primary transition-colors hover:text-primary/80">Sign in</Link>
              </p>
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
