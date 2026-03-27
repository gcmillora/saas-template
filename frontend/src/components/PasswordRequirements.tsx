import { cn } from "@/lib/utils";

interface PasswordRequirementsProps {
  password: string;
}

const requirements = [
  { label: "At least 8 characters", test: (p: string) => p.length >= 8 },
  { label: "1 uppercase letter", test: (p: string) => /[A-Z]/.test(p) },
  { label: "1 lowercase letter", test: (p: string) => /[a-z]/.test(p) },
  { label: "1 number", test: (p: string) => /\d/.test(p) },
  {
    label: "1 special character",
    test: (p: string) => /[^a-zA-Z0-9\s]/.test(p),
  },
];

export function PasswordRequirements({ password }: PasswordRequirementsProps) {
  if (!password) return null;

  return (
    <ul className="space-y-1 text-xs text-muted-foreground">
      {requirements.map((req) => {
        const met = req.test(password);
        return (
          <li
            key={req.label}
            className={cn(
              "flex items-center gap-1.5 transition-colors",
              met && "text-emerald-600"
            )}
          >
            <span className="text-[10px]">{met ? "\u2713" : "\u2022"}</span>
            {req.label}
          </li>
        );
      })}
    </ul>
  );
}
