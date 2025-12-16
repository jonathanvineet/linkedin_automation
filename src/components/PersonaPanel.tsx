import { motion } from "framer-motion";
import { User, Zap, MousePointer, AlertTriangle } from "lucide-react";
import { Slider } from "@/components/ui/slider";
import { cn } from "@/lib/utils";

interface PersonaPanelProps {
  selectedPersona: string;
  onPersonaChange: (persona: string) => void;
  settings: {
    typingSpeed: number;
    mousePrecision: number;
    errorRate: number;
  };
  onSettingsChange: (settings: any) => void;
}

const personas = [
  { id: "recruiter", label: "Recruiter", icon: User },
  { id: "founder", label: "Founder", icon: Zap },
  { id: "sales", label: "Sales", icon: MousePointer },
];

const behaviorSummary = {
  recruiter: "Methodical browsing patterns, extended profile views, professional messaging cadence.",
  founder: "Quick scanning, decisive actions, shorter session durations with high-value targeting.",
  sales: "High engagement rate, personalized outreach patterns, consistent follow-up timing.",
};

export function PersonaPanel({
  selectedPersona,
  onPersonaChange,
  settings,
  onSettingsChange,
}: PersonaPanelProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.1 }}
      className="rounded-lg border border-border bg-card p-5"
    >
      <div className="mb-4 flex items-center gap-2">
        <User className="h-4 w-4 text-primary" />
        <h3 className="font-mono text-sm font-semibold uppercase tracking-wider text-foreground">
          Persona Configuration
        </h3>
      </div>

      <div className="mb-6 flex gap-2">
        {personas.map((persona) => {
          const Icon = persona.icon;
          const isSelected = selectedPersona === persona.id;
          return (
            <button
              key={persona.id}
              onClick={() => onPersonaChange(persona.id)}
              className={cn(
                "flex flex-1 items-center justify-center gap-2 rounded-md border px-3 py-2 text-sm font-medium transition-all",
                isSelected
                  ? "border-primary bg-primary/10 text-primary"
                  : "border-border bg-secondary text-muted-foreground hover:border-muted-foreground hover:text-foreground"
              )}
            >
              <Icon className="h-4 w-4" />
              {persona.label}
            </button>
          );
        })}
      </div>

      <div className="space-y-5">
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <label className="text-xs font-medium text-muted-foreground">
              Typing Speed
            </label>
            <span className="font-mono text-xs text-primary">
              {settings.typingSpeed} WPM
            </span>
          </div>
          <Slider
            value={[settings.typingSpeed]}
            onValueChange={([value]) =>
              onSettingsChange({ ...settings, typingSpeed: value })
            }
            min={20}
            max={120}
            step={1}
            className="cursor-pointer"
          />
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <label className="text-xs font-medium text-muted-foreground">
              Mouse Precision
            </label>
            <span className="font-mono text-xs text-primary">
              {settings.mousePrecision}%
            </span>
          </div>
          <Slider
            value={[settings.mousePrecision]}
            onValueChange={([value]) =>
              onSettingsChange({ ...settings, mousePrecision: value })
            }
            min={60}
            max={100}
            step={1}
            className="cursor-pointer"
          />
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <label className="text-xs font-medium text-muted-foreground">
              Error Rate
            </label>
            <span className="font-mono text-xs text-warning">
              {settings.errorRate}%
            </span>
          </div>
          <Slider
            value={[settings.errorRate]}
            onValueChange={([value]) =>
              onSettingsChange({ ...settings, errorRate: value })
            }
            min={0}
            max={15}
            step={0.5}
            className="cursor-pointer"
          />
        </div>
      </div>

      <div className="mt-5 rounded-md border border-border bg-secondary/50 p-3">
        <div className="mb-2 flex items-center gap-2">
          <AlertTriangle className="h-3 w-3 text-warning" />
          <span className="text-xs font-medium text-muted-foreground">
            Human Behavior Summary
          </span>
        </div>
        <p className="font-mono text-xs leading-relaxed text-foreground/80">
          {behaviorSummary[selectedPersona as keyof typeof behaviorSummary]}
        </p>
      </div>
    </motion.div>
  );
}
