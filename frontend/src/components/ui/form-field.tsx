import * as React from "react";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

interface FormFieldProps {
  label: string;
  required?: boolean;
  description?: string;
  error?: string;
  className?: string;
  children: React.ReactNode;
  htmlFor?: string;
}

function FormField({
  label,
  required,
  description,
  error,
  className,
  children,
  htmlFor,
}: FormFieldProps) {
  return (
    <div className={cn("space-y-1.5", className)}>
      <div className="flex items-center gap-1.5">
        <Label htmlFor={htmlFor} className="text-sm font-medium">
          {label}
        </Label>
        {required && (
          <span className="size-1.5 rounded-full bg-primary/70" aria-hidden />
        )}
      </div>
      {children}
      {description && !error && (
        <p className="text-xs text-muted-foreground/80 leading-normal">
          {description}
        </p>
      )}
      {error && (
        <p className="text-xs text-destructive leading-normal">{error}</p>
      )}
    </div>
  );
}

interface FormSectionProps {
  title?: string;
  children: React.ReactNode;
  className?: string;
}

function FormSection({ title, children, className }: FormSectionProps) {
  return (
    <div className={cn("space-y-4", className)}>
      {title && (
        <div className="flex items-center gap-3">
          <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground/60">
            {title}
          </span>
          <div className="h-px flex-1 bg-border/60" />
        </div>
      )}
      {children}
    </div>
  );
}

export { FormField, FormSection };
