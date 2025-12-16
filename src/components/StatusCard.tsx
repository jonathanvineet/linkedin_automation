import { motion } from "framer-motion";
import { LucideIcon } from "lucide-react";
import { cn } from "@/lib/utils";

interface StatusCardProps {
  title: string;
  value: string | number;
  icon: LucideIcon;
  status: "active" | "inactive" | "warning" | "primary";
  subtitle?: string;
}

const statusStyles = {
  active: "border-success/30 status-glow-active",
  inactive: "border-destructive/30 status-glow-inactive",
  warning: "border-warning/30 status-glow-warning",
  primary: "border-primary/30 status-glow-primary",
};

const indicatorStyles = {
  active: "bg-success",
  inactive: "bg-destructive",
  warning: "bg-warning",
  primary: "bg-primary",
};

export function StatusCard({ title, value, icon: Icon, status, subtitle }: StatusCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      className={cn(
        "relative overflow-hidden rounded-lg border bg-card p-4 transition-all duration-300",
        statusStyles[status]
      )}
    >
      <div className="flex items-start justify-between">
        <div className="space-y-1">
          <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground">
            {title}
          </p>
          <p className="font-mono text-2xl font-semibold text-foreground">{value}</p>
          {subtitle && (
            <p className="text-xs text-muted-foreground">{subtitle}</p>
          )}
        </div>
        <div className="flex items-center gap-2">
          <span
            className={cn(
              "h-2 w-2 rounded-full animate-pulse-glow",
              indicatorStyles[status]
            )}
          />
          <Icon className="h-5 w-5 text-muted-foreground" />
        </div>
      </div>
    </motion.div>
  );
}
