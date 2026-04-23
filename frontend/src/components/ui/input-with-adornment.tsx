import * as React from "react";
import { cn } from "@/lib/utils";

interface InputWithAdornmentProps
  extends Omit<React.ComponentProps<"input">, "prefix"> {
  prefix?: React.ReactNode;
  suffix?: React.ReactNode;
}

function InputWithAdornment({
  className,
  prefix,
  suffix,
  type,
  ...props
}: InputWithAdornmentProps) {
  return (
    <div
      className={cn(
        "flex items-center h-9 w-full rounded-md border border-input bg-transparent shadow-xs transition-[color,box-shadow] outline-none",
        "has-[:focus-visible]:border-ring has-[:focus-visible]:ring-ring/50 has-[:focus-visible]:ring-[3px]",
        "has-[:disabled]:opacity-50 has-[:disabled]:cursor-not-allowed",
        "has-aria-invalid:ring-destructive/20 has-aria-invalid:border-destructive",
      )}
    >
      {prefix && (
        <span className="flex items-center pl-3 text-sm text-muted-foreground/70 select-none">
          {prefix}
        </span>
      )}
      <input
        type={type}
        className={cn(
          "flex-1 h-full bg-transparent px-3 py-1 text-base outline-none placeholder:text-muted-foreground md:text-sm min-w-0",
          "file:text-foreground file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium",
          "disabled:pointer-events-none disabled:cursor-not-allowed",
          prefix && "pl-1.5",
          suffix && "pr-1.5",
          className,
        )}
        {...props}
      />
      {suffix && (
        <span className="flex items-center pr-3 text-sm text-muted-foreground/70 select-none">
          {suffix}
        </span>
      )}
    </div>
  );
}

export { InputWithAdornment };
