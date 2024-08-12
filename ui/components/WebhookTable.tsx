import React, { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { WebhookDetailsDialog } from "@/components/WebhookDetailsDialog";
import { useWebhookLogs } from "@/hooks/useWebhookLogs";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export const WebhookTable: React.FC = () => {
  const { logs, loading, error, fetchNextPage } = useWebhookLogs();
  const [selectedLogId, setSelectedLogId] = useState<string | null>(null);

  if (loading) return <TableSkeleton />;
  if (error) return <ErrorDisplay message={error.message} />;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Webhook Logs</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Timestamp</TableHead>
              <TableHead>Method</TableHead>
              <TableHead>Path</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Duration (ms)</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {logs.map((log) => (
              <TableRow
                key={log.id}
                onClick={() => setSelectedLogId(log.id)}
                className="cursor-pointer hover:bg-muted"
              >
                <TableCell>
                  {new Date(log.timestamp).toLocaleString()}
                </TableCell>
                <TableCell>{log.method}</TableCell>
                <TableCell>{log.path}</TableCell>
                <TableCell>{log.status_code}</TableCell>
                <TableCell>{log.duration}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <div className="mt-4 flex justify-center">
          <Button onClick={fetchNextPage}>Load More</Button>
        </div>
      </CardContent>
      {selectedLogId && (
        <WebhookDetailsDialog
          logId={selectedLogId}
          onClose={() => setSelectedLogId(null)}
        />
      )}
    </Card>
  );
};

const TableSkeleton: React.FC = () => (
  <Card>
    <CardHeader>
      <CardTitle>Webhook Logs</CardTitle>
    </CardHeader>
    <CardContent>
      <div className="space-y-2">
        <div className="h-8 bg-gray-200 rounded animate-pulse" />
        <div className="h-8 bg-gray-200 rounded animate-pulse" />
        <div className="h-8 bg-gray-200 rounded animate-pulse" />
        <div className="h-8 bg-gray-200 rounded animate-pulse" />
        <div className="h-8 bg-gray-200 rounded animate-pulse" />
      </div>
    </CardContent>
  </Card>
);

const ErrorDisplay: React.FC<{ message: string }> = ({ message }) => (
  <Card>
    <CardHeader>
      <CardTitle>Error</CardTitle>
    </CardHeader>
    <CardContent>
      <p className="text-red-500">{message}</p>
    </CardContent>
  </Card>
);
