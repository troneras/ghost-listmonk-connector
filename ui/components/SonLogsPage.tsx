import React, { useState } from "react";
import { useSonLogs } from "@/hooks/useSonLogs";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
  PaginationEllipsis,
} from "@/components/ui/pagination";
import { ActionExecutionLog } from "@/lib/types";
import { useToast } from "@/components/ui/use-toast";
import Link from "next/link";
import { WebhookDetailsDialog } from "@/components/WebhookDetailsDialog";
import { Badge } from "@/components/ui/badge";

const StatusBadge: React.FC<{ status: string }> = ({ status }) => {
  let color = "bg-gray-500";
  switch (status.toLowerCase()) {
    case "success":
      color = "bg-green-500";
      break;
    case "failure":
      color = "bg-red-500";
      break;
    case "warning":
      color = "bg-yellow-500";
      break;
    case "pending":
      color = "bg-blue-500";
      break;
  }
  return <Badge className={`${color} text-white`}>{status}</Badge>;
};

const SonLogsPage: React.FC = () => {
  const { logs, loading, error, pagination, fetchLogs, fetchActionLogs } =
    useSonLogs();
  const [selectedExecution, setSelectedExecution] = useState<string | null>(
    null
  );
  const [actionLogs, setActionLogs] = useState<ActionExecutionLog[]>([]);
  const [selectedWebhookLogId, setSelectedWebhookLogId] = useState<
    string | null
  >(null);
  const { toast } = useToast();

  const totalPages = Math.ceil(pagination.total / pagination.limit);
  const currentPage = Math.floor(pagination.offset / pagination.limit) + 1;

  const handlePageChange = (page: number) => {
    const newOffset = (page - 1) * pagination.limit;
    fetchLogs(newOffset);
  };

  const openActionLogs = async (executionId: string) => {
    try {
      setSelectedExecution(executionId);
      const actionLogsData = await fetchActionLogs(executionId);
      setActionLogs(actionLogsData);
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to fetch action logs",
        variant: "destructive",
      });
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  return (
    <div className="container mx-auto py-10">
      <Card>
        <CardHeader>
          <CardTitle>Son Execution Logs</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Son Name</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Executed At</TableHead>
                <TableHead>Webhook</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {logs.map((log) => (
                <TableRow key={log.id}>
                  <TableCell>
                    <Link
                      href={`/sons/${log.son_id}`}
                      className="text-blue-600 hover:underline"
                    >
                      {log.sonName}
                    </Link>
                  </TableCell>
                  <TableCell>
                    <StatusBadge status={log.status} />
                  </TableCell>
                  <TableCell>
                    {new Date(log.executed_at).toLocaleString()}
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="link"
                      onClick={() =>
                        setSelectedWebhookLogId(log.webhook_log_id)
                      }
                    >
                      View Webhook
                    </Button>
                  </TableCell>
                  <TableCell>
                    <Button onClick={() => openActionLogs(log.id)}>
                      View Actions
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          <Pagination>
            <PaginationContent>
              {currentPage > 1 && (
                <PaginationItem>
                  <PaginationPrevious
                    onClick={() => handlePageChange(currentPage - 1)}
                  />
                </PaginationItem>
              )}
              {[...Array(totalPages)].map((_, index) => {
                const page = index + 1;
                if (
                  page === 1 ||
                  page === totalPages ||
                  (page >= currentPage - 1 && page <= currentPage + 1)
                ) {
                  return (
                    <PaginationItem key={page}>
                      <PaginationLink
                        isActive={page === currentPage}
                        onClick={() => handlePageChange(page)}
                      >
                        {page}
                      </PaginationLink>
                    </PaginationItem>
                  );
                } else if (
                  page === currentPage - 2 ||
                  page === currentPage + 2
                ) {
                  return <PaginationEllipsis key={page} />;
                }
                return null;
              })}
              {currentPage < totalPages && (
                <PaginationItem>
                  <PaginationNext
                    onClick={() => handlePageChange(currentPage + 1)}
                  />
                </PaginationItem>
              )}
            </PaginationContent>
          </Pagination>
        </CardContent>
      </Card>

      <Dialog
        open={!!selectedExecution}
        onOpenChange={() => setSelectedExecution(null)}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              Action Logs for Execution {selectedExecution}
            </DialogTitle>
          </DialogHeader>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Action Type</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Executed At</TableHead>
                <TableHead>Error Message</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {actionLogs.map((log: ActionExecutionLog) => (
                <TableRow key={log.id}>
                  <TableCell>{log.action_type}</TableCell>
                  <TableCell>
                    <StatusBadge status={log.status} />
                  </TableCell>
                  <TableCell>
                    {new Date(log.executed_at).toLocaleString()}
                  </TableCell>
                  <TableCell>{log.error_message || "N/A"}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </DialogContent>
      </Dialog>

      {selectedWebhookLogId && (
        <WebhookDetailsDialog
          logId={selectedWebhookLogId}
          onClose={() => setSelectedWebhookLogId(null)}
        />
      )}
    </div>
  );
};

export default SonLogsPage;
