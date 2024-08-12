// ui/components/Dashboard.tsx

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from "recharts";
import { ChartContainer, ChartConfig } from "@/components/ui/chart";
import { DashboardSkeleton } from "./DashboardSKeleton";
import { useRecentActivity } from "@/hooks/useRecentActivity";
import { useSonStats } from "@/hooks/useSonStats";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const chartConfig = {
  executions: {
    label: "Executions",
    color: "#2563eb",
  },
  success: {
    label: "Success",
    color: "#F97316",
  },
  failure: {
    label: "Failure",
    color: "#d1d5db",
  },
} satisfies ChartConfig;

const timeframeOptions = [
  { value: "1h", label: "1 Hour" },
  { value: "6h", label: "6 Hours" },
  { value: "12h", label: "12 Hours" },
  { value: "24h", label: "24 Hours" },
  { value: "168h", label: "1 Week" },
];

export default function Dashboard() {
  const [timeframe, setTimeframe] = useState("24h");
  const {
    activities,
    loading: activitiesLoading,
    error: activitiesError,
  } = useRecentActivity();
  const {
    stats,
    loading: statsLoading,
    error: statsError,
  } = useSonStats(timeframe);

  if (activitiesLoading || statsLoading) {
    return <DashboardSkeleton />;
  }

  if (activitiesError || statsError) {
    return <div>Error loading dashboard data</div>;
  }

  return (
    <div className="flex-1 space-y-4">
      <h1 className="text-2xl font-bold">Dashboard</h1>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
          </CardHeader>
          <CardContent>
            {!activities && <div>No recent activity</div>}
            {activities &&
              activities.map((activity) => (
                <Alert key={activity.id} className="mb-2">
                  <AlertTitle>{activity.action_type}</AlertTitle>
                  <AlertDescription>
                    {activity.description}
                    <br />
                    {new Date(activity.timestamp).toLocaleString()}
                  </AlertDescription>
                </Alert>
              ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle>Son Performance</CardTitle>
            <Select value={timeframe} onValueChange={setTimeframe}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Select timeframe" />
              </SelectTrigger>
              <SelectContent>
                {timeframeOptions.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </CardHeader>
          <CardContent>
            <ChartContainer
              config={chartConfig}
              className="min-h-[200px] w-full"
            >
              <BarChart data={stats}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="success" fill="var(--color-success)" />
                <Bar dataKey="failure" fill="var(--color-failure)" />
              </BarChart>
            </ChartContainer>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
