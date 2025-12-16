import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { 
  Globe, 
  UserCheck, 
  Bot, 
  Shield, 
  Send, 
  MessageSquare, 
  Clock,
  AlertTriangle,
  Terminal
} from "lucide-react";
import { StatusCard } from "@/components/StatusCard";
import { PersonaPanel } from "@/components/PersonaPanel";
import { AutomationFlow } from "@/components/AutomationFlow";
import { StealthTechniques } from "@/components/StealthTechniques";
import { ActivityLog } from "@/components/ActivityLog";

const Dashboard = () => {
  const [currentStep, setCurrentStep] = useState(3);
  const [selectedPersona, setSelectedPersona] = useState("recruiter");
  const [settings, setSettings] = useState({
    typingSpeed: 65,
    mousePrecision: 87,
    errorRate: 3.5,
  });

  // Simulate step progression
  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentStep((prev) => (prev >= 5 ? 1 : prev + 1));
    }, 4000);
    return () => clearInterval(interval);
  }, []);

  const stealthTechniques = [
    { id: "mouse-curves", label: "Mouse Curves", description: "Bézier curve interpolation", icon: Globe, enabled: true },
    { id: "typing-variance", label: "Typing Variance", description: "Random keystroke delays", icon: Globe, enabled: true },
    { id: "scroll-randomization", label: "Scroll Randomization", description: "Non-linear scroll patterns", icon: Globe, enabled: true },
    { id: "fingerprint-masking", label: "Fingerprint Masking", description: "Canvas & WebGL noise", icon: Globe, enabled: false },
  ];

  return (
    <div className="min-h-screen bg-background noise-overlay">
      {/* Warning Banner */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="border-b border-warning/30 bg-warning/10 px-4 py-2"
      >
        <div className="mx-auto flex max-w-7xl items-center justify-center gap-2">
          <AlertTriangle className="h-4 w-4 text-warning" />
          <p className="font-mono text-xs font-medium text-warning">
            EDUCATIONAL DEMO — DO NOT USE IN PRODUCTION
          </p>
        </div>
      </motion.div>

      {/* Header */}
      <header className="border-b border-border bg-card/50 backdrop-blur-sm">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-lg border border-primary/30 bg-primary/10">
              <Terminal className="h-5 w-5 text-primary" />
            </div>
            <div>
              <h1 className="font-mono text-lg font-semibold text-foreground">
                AutoPilot<span className="text-primary">.dev</span>
              </h1>
              <p className="text-xs text-muted-foreground">
                Automation Control Interface v0.1.0
              </p>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 rounded-md border border-border bg-secondary px-3 py-1.5">
              <span className="h-2 w-2 animate-pulse-glow rounded-full bg-success" />
              <span className="font-mono text-xs text-muted-foreground">
                System Active
              </span>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="mx-auto max-w-7xl px-6 py-6">
        {/* Status Cards */}
        <div className="mb-6 grid grid-cols-2 gap-4 lg:grid-cols-4">
          <StatusCard
            title="Browser Session"
            value="Active"
            icon={Globe}
            status="active"
            subtitle="Chromium 120.0"
          />
          <StatusCard
            title="Logged In"
            value="Yes"
            icon={UserCheck}
            status="active"
            subtitle="Session valid"
          />
          <StatusCard
            title="Persona Active"
            value={selectedPersona.charAt(0).toUpperCase() + selectedPersona.slice(1)}
            icon={Bot}
            status="primary"
            subtitle="Behavior loaded"
          />
          <StatusCard
            title="Stealth Mode"
            value="Enabled"
            icon={Shield}
            status="active"
            subtitle="3/4 techniques"
          />
        </div>

        {/* Daily Counters */}
        <div className="mb-6 grid grid-cols-3 gap-4">
          <StatusCard
            title="Connections Sent"
            value="47"
            icon={Send}
            status="primary"
            subtitle="Daily limit: 100"
          />
          <StatusCard
            title="Messages Sent"
            value="23"
            icon={MessageSquare}
            status="primary"
            subtitle="Daily limit: 50"
          />
          <StatusCard
            title="Cooldown Time"
            value="2m 34s"
            icon={Clock}
            status="warning"
            subtitle="Next action at 14:38"
          />
        </div>

        {/* Main Grid */}
        <div className="grid gap-6 lg:grid-cols-2">
          {/* Left Column */}
          <div className="space-y-6">
            <PersonaPanel
              selectedPersona={selectedPersona}
              onPersonaChange={setSelectedPersona}
              settings={settings}
              onSettingsChange={setSettings}
            />
            <StealthTechniques techniques={stealthTechniques} />
          </div>

          {/* Right Column */}
          <div className="space-y-6">
            <AutomationFlow currentStep={currentStep} />
            <ActivityLog entries={undefined as any} />
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="mt-8 border-t border-border bg-card/30">
        <div className="mx-auto max-w-7xl px-6 py-4">
          <p className="text-center font-mono text-xs text-muted-foreground">
            Internal Development Tool — All data is mock/simulated — No real automation occurs
          </p>
        </div>
      </footer>
    </div>
  );
};

export default Dashboard;
