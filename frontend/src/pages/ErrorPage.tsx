import { useRouteError, isRouteErrorResponse, Link } from "react-router";
import { Button } from "@/components/ui/button";

export function ErrorPage() {
  const error = useRouteError();
  const isNotFound = isRouteErrorResponse(error) && error.status === 404;

  if (isNotFound) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center px-6">
        <span className="text-[10rem] leading-none font-black tracking-tighter text-foreground/10">404</span>
        <h1 className="mt-2 text-2xl font-bold tracking-tight text-foreground">Page not found</h1>
        <p className="mt-2 text-sm text-muted-foreground">The page you're looking for doesn't exist or has been moved.</p>
        <Button asChild className="mt-8"><Link to="/">Go Home</Link></Button>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center px-6">
      <h1 className="mt-6 text-2xl font-bold tracking-tight text-foreground">Something went wrong</h1>
      <p className="mt-2 text-sm text-muted-foreground">An unexpected error occurred. Please try again.</p>
      <div className="mt-8 flex gap-3">
        <Button variant="outline" onClick={() => window.location.reload()}>Refresh</Button>
        <Button asChild><Link to="/">Go Home</Link></Button>
      </div>
    </div>
  );
}
