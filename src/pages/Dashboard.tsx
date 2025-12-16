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
  Terminal,
  Play,
  Square
} from "lucide-react";
import { StatusCard } from "@/components/StatusCard";
import { PersonaPanel } from "@/components/PersonaPanel";
import { AutomationFlow } from "@/components/AutomationFlow";
import { StealthTechniques } from "@/components/StealthTechniques";
import { ActivityLog } from "@/components/ActivityLog";
import { apiClient, type SystemStatus, type Stats, type ActivityLogEntry } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { useToast } from "@/hooks/use-toast";

const Dashboard = () => {
  const [currentStep, setCurrentStep] = useState(3);
  const [selectedPersona, setSelectedPersona] = useState("recruiter");
  const [settings, setSettings] = useState({
    typingSpeed: 65,
    mousePrecision: 87,
    errorRate: 3.5,
  });
  const [systemStatus, setSystemStatus] = useState<SystemStatus | null>(null);
  const [stats, setStats] = useState<Stats | null>(null);
  const [activityLogs, setActivityLogs] = useState<ActivityLogEntry[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const { toast } = useToast();

  // Fetch status and stats periodically
  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000); // Every 5 seconds
    return () => clearInterval(interval);
  }, []);

  const fetchData = async () => {
    try {
      const [status, statistics, logs] = await Promise.all([
        apiClient.getStatus(),
        apiClient.getStats(),
        apiClient.getActivity(),
      ]);
      
      setSystemStatus(status);
      setStats(statistics);
      setActivityLogs(logs);
    } catch (error) {
      console.error("Failed to fetch data:", error);
    }
  };

  const handleStart = async () => {
    setIsLoading(true);
    try {
      const response = await apiClient.start();
      toast({
        title: "Success",
        description: response.message || "Automation started successfully",
      });
      fetchData();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to start automation",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleStop = async () => {
    setIsLoading(true);
    try {
      const response = await apiClient.stop();
      toast({
        title: "Success",
        description: response.message || "Automation stopped successfully",
      });
      fetchData();
    } catch (error: any) {
      toast({
        title: "Error",
        description: error.message || "Failed to stop automation",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handlePersonaChange = async (persona: string) => {
    setSelectedPersona(persona);
    try {
      await apiClient.setPersona(persona);
      fetchData();
    } catch (error) {
      console.error("Failed to change persona:", error);
    }
  };

  // Format cooldown time
  const formatCooldown = (seconds: number): string => {
    if (seconds <= 0) return "Ready";
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}m ${secs}s`;
  };

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
              <span className={`h-2 w-2 animate-pulse-glow rounded-full ${systemStatus?.running ? 'bg-success' : 'bg-muted-foreground'}`} />
              <span className="font-mono text-xs text-muted-foreground">
                {systemStatus?.running ? 'System Active' : 'System Idle'}
              </span>
            </div>
            
            {systemStatus?.running ? (
              <Button
                onClick={handleStop}
                disabled={isLoading}
                variant="destructive"
                size="sm"
              >
                <Square className="h-4 w-4 mr-2" />
                Stop
              </Button>
            ) : (
              <Button
                onClick={handleStart}
                disabled={isLoading}
                variant="default"
                size="sm"
              >
                <Play className="h-4 w-4 mr-2" />
                Start
              </Button>
            )}
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="mx-auto max-w-7xl px-6 py-6">
        {/* Status Cards */}
        <div className="mb-6 grid grid-cols-2 gap-4 lg:grid-cols-4">
          <StatusCard
            title="Browser Session"
            value={systemStatus?.running ? "Active" : "Idle"}
            icon={Globe}
            status={systemStatus?.running ? "active" : "inactive"}
            subtitle="Chromium 120.0"
          />
          <StatusCard
            title="Logged In"
            value={systemStatus?.logged_in ? "Yes" : "No"}
            icon={UserCheck}
            status={systemStatus?.logged_in ? "active" : "inactive"}
            subtitle={systemStatus?.logged_in ? "Session valid" : "Not connected"}
          />
          <StatusCard
            title="Persona Active"
            value={systemStatus?.persona || selectedPersona.charAt(0).toUpperCase() + selectedPersona.slice(1)}
            icon={Bot}
            status="primary"
            subtitle="Behavior loaded"
          />
          <StatusCard
            title="Stealth Mode"
            value={systemStatus?.stealth ? "Enabled" : "Disabled"}
            icon={Shield}
            status={systemStatus?.stealth ? "active" : "inactive"}
            subtitle="8 techniques"
          />
        </div>

        {/* Daily Counters */}
        <div className="mb-6 grid grid-cols-3 gap-4">
          <StatusCard
            title="Connections Sent"
            value={stats?.connections_sent || 0}
            icon={Send}
            status="primary"
            subtitle={`Daily limit: ${stats?.daily_limit?.connections || 20}`}
          />
          <StatusCard
            title="Messages Sent"
            value={stats?.messages_sent || 0}
            icon={MessageSquare}
            status="primary"
            subtitle={`Daily limit: ${stats?.daily_limit?.messages || 10}`}
          />
          <StatusCard
            title="Cooldown Time"
            value={formatCooldown(stats?.cooldown_seconds || 0)}
            icon={Clock}
            status={stats && stats.cooldown_seconds > 0 ? "warning" : "active"}
            subtitle={stats && stats.cooldown_seconds > 0 ? "Action in cooldown" : "Ready to act"}
          />
        </div>

        {/* Main Grid */}
        <div className="grid gap-6 lg:grid-cols-2">
          {/* Left Column */}
          <div className="space-y-6">
            <PersonaPanel
              selectedPersona={selectedPersona}
              onPersonaChange={handlePersonaChange}
              settings={settings}
              onSettingsChange={setSettings}
            />
            <StealthTechniques techniques={stealthTechniques} />
          </div>

          {/* Right Column */}
          <div className="space-y-6">
            <AutomationFlow currentStep={currentStep} />
            <ActivityLog entries={activityLogs} />
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
