import { motion } from "framer-motion";
import { Activity, LogIn, Search, UserCircle, UserPlus, MessageSquare, Clock, AlertCircle } from "lucide-react";
import { cn } from "@/lib/utils";

interface LogEntry {
  id: string;
  timestamp: string;
  action: string;
  type: "info" | "success" | "warning" | "error";
  details?: string;
}

interface ActivityLogProps {
  entries: LogEntry[];
}

const typeIcons = {
  info: Clock,
  success: UserPlus,
  warning: AlertCircle,
  error: AlertCircle,
};

const typeStyles = {
  info: "text-muted-foreground border-border",
  success: "text-success border-success/30",
  warning: "text-warning border-warning/30",
  error: "text-destructive border-destructive/30",
};

const mockEntries: LogEntry[] = [
  { id: "1", timestamp: "14:32:15", action: "Session initialized", type: "info", details: "Browser context created" },
  { id: "2", timestamp: "14:32:18", action: "Login successful", type: "success", details: "2FA bypassed via cookie" },
  { id: "3", timestamp: "14:33:02", action: "Search executed", type: "info", details: "\"Software Engineer\" in Bay Area" },
  { id: "4", timestamp: "14:33:45", action: "Profile visited", type: "info", details: "John D. — Senior Developer" },
  { id: "5", timestamp: "14:34:12", action: "Connection sent", type: "success", details: "Personalized note attached" },
  { id: "6", timestamp: "14:35:01", action: "Rate limit warning", type: "warning", details: "Approaching daily limit" },
  { id: "7", timestamp: "14:35:30", action: "Cooldown initiated", type: "info", details: "Waiting 45 seconds" },
  { id: "8", timestamp: "14:36:15", action: "Profile visited", type: "info", details: "Sarah M. — Tech Lead" },
];

export function ActivityLog({ entries = mockEntries }: ActivityLogProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.4 }}
      className="rounded-lg border border-border bg-card p-5"
    >
      <div className="mb-4 flex items-center gap-2">
        <Activity className="h-4 w-4 text-primary" />
        <h3 className="font-mono text-sm font-semibold uppercase tracking-wider text-foreground">
          Activity Log
        </h3>
        <span className="ml-auto font-mono text-xs text-muted-foreground">
          {entries.length} entries
        </span>
      </div>

      <div className="max-h-[320px] space-y-1 overflow-y-auto pr-2 scrollbar-thin">
        {entries.map((entry, index) => {
          const Icon = typeIcons[entry.type];
          return (
            <motion.div
              key={entry.id}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: index * 0.05 }}
              className={cn(
                "flex items-start gap-3 rounded-md border-l-2 bg-secondary/30 px-3 py-2",
                typeStyles[entry.type]
              )}
            >
              <span className="mt-0.5 font-mono text-xs text-muted-foreground">
                {entry.timestamp}
              </span>
              
              <Icon className={cn("mt-0.5 h-3.5 w-3.5 shrink-0", typeStyles[entry.type].split(" ")[0])} />
              
              <div className="min-w-0 flex-1">
                <p className="font-mono text-xs font-medium text-foreground">
                  {entry.action}
                </p>
                {entry.details && (
                  <p className="truncate text-xs text-muted-foreground">
                    {entry.details}
                  </p>
                )}
              </div>
            </motion.div>
          );
        })}
      </div>
    </motion.div>
  );
}
