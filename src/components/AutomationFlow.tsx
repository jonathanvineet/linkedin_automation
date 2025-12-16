import { motion } from "framer-motion";
import { LogIn, Search, UserCircle, UserPlus, MessageSquare, CheckCircle } from "lucide-react";
import { cn } from "@/lib/utils";

interface AutomationFlowProps {
  currentStep: number;
}

const steps = [
  { id: 1, label: "Login", icon: LogIn, description: "Initialize session" },
  { id: 2, label: "Search", icon: Search, description: "Query targets" },
  { id: 3, label: "Visit Profile", icon: UserCircle, description: "View candidate" },
  { id: 4, label: "Connect", icon: UserPlus, description: "Send request" },
  { id: 5, label: "Message", icon: MessageSquare, description: "Outreach" },
];

export function AutomationFlow({ currentStep }: AutomationFlowProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.2 }}
      className="rounded-lg border border-border bg-card p-5"
    >
      <div className="mb-4 flex items-center gap-2">
        <CheckCircle className="h-4 w-4 text-primary" />
        <h3 className="font-mono text-sm font-semibold uppercase tracking-wider text-foreground">
          Automation Flow
        </h3>
      </div>

      <div className="space-y-1">
        {steps.map((step, index) => {
          const Icon = step.icon;
          const isActive = step.id === currentStep;
          const isComplete = step.id < currentStep;
          const isPending = step.id > currentStep;

          return (
            <div key={step.id} className="relative">
              <div
                className={cn(
                  "flex items-center gap-3 rounded-md px-3 py-2.5 transition-all",
                  isActive && "bg-primary/10 border border-primary/30",
                  isComplete && "opacity-60",
                  isPending && "opacity-40"
                )}
              >
                <div
                  className={cn(
                    "flex h-8 w-8 items-center justify-center rounded-md border",
                    isActive && "border-primary bg-primary/20 text-primary",
                    isComplete && "border-success bg-success/20 text-success",
                    isPending && "border-border bg-secondary text-muted-foreground"
                  )}
                >
                  {isComplete ? (
                    <CheckCircle className="h-4 w-4" />
                  ) : (
                    <Icon className="h-4 w-4" />
                  )}
                </div>

                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <span
                      className={cn(
                        "font-mono text-sm font-medium",
                        isActive && "text-primary",
                        isComplete && "text-success",
                        isPending && "text-muted-foreground"
                      )}
                    >
                      {step.label}
                    </span>
                    {isActive && (
                      <motion.span
                        initial={{ opacity: 0, scale: 0.8 }}
                        animate={{ opacity: 1, scale: 1 }}
                        className="rounded-full bg-primary px-2 py-0.5 text-[10px] font-semibold uppercase text-primary-foreground"
                      >
                        Active
                      </motion.span>
                    )}
                  </div>
                  <p className="text-xs text-muted-foreground">
                    {step.description}
                  </p>
                </div>

                <span className="font-mono text-xs text-muted-foreground">
                  {String(step.id).padStart(2, "0")}
                </span>
              </div>

              {index < steps.length - 1 && (
                <div
                  className={cn(
                    "ml-6 h-2 w-px",
                    isComplete ? "bg-success/40" : "bg-border"
                  )}
                />
              )}
            </div>
          );
        })}
      </div>
    </motion.div>
  );
}
