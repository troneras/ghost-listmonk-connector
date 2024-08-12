"use client";
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

// Mock data - replace with actual API calls in production
const recentActivity = [
  {
    id: 1,
    action: 'Son "Welcome Email" triggered',
    timestamp: "2024-08-10T10:30:00Z",
  },
  {
    id: 2,
    action: 'New Son "Survey Request" created',
    timestamp: "2024-08-10T09:15:00Z",
  },
  {
    id: 3,
    action: 'Son "Monthly Newsletter" modified',
    timestamp: "2024-08-09T16:45:00Z",
  },
];

const sonStats = [
  { name: "Welcome Email", executions: 120, success: 115, failure: 5 },
  { name: "Survey Request", executions: 80, success: 78, failure: 2 },
  { name: "Monthly Newsletter", executions: 50, success: 50, failure: 0 },
];

const chartConfig = {
  executions: {
    label: "Executions",
    color: "#2563eb",
  },
  success: {
    label: "Success",
    color: "#60a5fa",
  },
  failure: {
    label: "Failure",
    color: "#d1d5db",
  },
} satisfies ChartConfig;

export default function Dashboard() {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Simulate API call
    setTimeout(() => setIsLoading(false), 1000);
  }, []);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Dashboard</h1>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
          </CardHeader>
          <CardContent>
            {recentActivity.map((activity) => (
              <Alert key={activity.id} className="mb-2">
                <AlertTitle>{activity.action}</AlertTitle>
                <AlertDescription>
                  {new Date(activity.timestamp).toLocaleString()}
                </AlertDescription>
              </Alert>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Son Performance</CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer
              config={chartConfig}
              className="min-h-[200px] w-full"
            >
              <BarChart data={sonStats}>
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
