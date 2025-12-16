import { motion } from "framer-motion";
import { Shield, MousePointer, Keyboard, ScrollText, Fingerprint } from "lucide-react";
import { Switch } from "@/components/ui/switch";
import { cn } from "@/lib/utils";

interface StealthTechnique {
  id: string;
  label: string;
  description: string;
  icon: any;
  enabled: boolean;
}

interface StealthTechniquesProps {
  techniques: StealthTechnique[];
}

const defaultTechniques: StealthTechnique[] = [
  {
    id: "mouse-curves",
    label: "Mouse Curves",
    description: "Bézier curve interpolation for natural movement",
    icon: MousePointer,
    enabled: true,
  },
  {
    id: "typing-variance",
    label: "Typing Variance",
    description: "Random delays between keystrokes",
    icon: Keyboard,
    enabled: true,
  },
  {
    id: "scroll-randomization",
    label: "Scroll Randomization",
    description: "Non-linear scroll patterns with pauses",
    icon: ScrollText,
    enabled: true,
  },
  {
    id: "fingerprint-masking",
    label: "Fingerprint Masking",
    description: "Canvas & WebGL noise injection",
    icon: Fingerprint,
    enabled: false,
  },
];

export function StealthTechniques({ techniques = defaultTechniques }: StealthTechniquesProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.3 }}
      className="rounded-lg border border-border bg-card p-5"
    >
      <div className="mb-4 flex items-center gap-2">
        <Shield className="h-4 w-4 text-primary" />
        <h3 className="font-mono text-sm font-semibold uppercase tracking-wider text-foreground">
          Stealth Techniques
        </h3>
        <span className="ml-auto rounded-full bg-secondary px-2 py-0.5 text-[10px] font-medium text-muted-foreground">
          Read-only
        </span>
      </div>

      <div className="space-y-3">
        {techniques.map((technique) => {
          const Icon = technique.icon;
          return (
            <div
              key={technique.id}
              className={cn(
                "flex items-center gap-3 rounded-md border px-3 py-2.5 transition-all",
                technique.enabled
                  ? "border-primary/20 bg-primary/5"
                  : "border-border bg-secondary/30"
              )}
            >
              <div
                className={cn(
                  "flex h-8 w-8 items-center justify-center rounded-md",
                  technique.enabled
                    ? "bg-primary/20 text-primary"
                    : "bg-secondary text-muted-foreground"
                )}
              >
                <Icon className="h-4 w-4" />
              </div>

              <div className="flex-1 min-w-0">
                <p
                  className={cn(
                    "font-mono text-sm font-medium",
                    technique.enabled ? "text-foreground" : "text-muted-foreground"
                  )}
                >
                  {technique.label}
                </p>
                <p className="truncate text-xs text-muted-foreground">
                  {technique.description}
                </p>
              </div>

              <Switch
                checked={technique.enabled}
                disabled
                className="pointer-events-none"
              />
            </div>
          );
        })}
      </div>
    </motion.div>
  );
}
