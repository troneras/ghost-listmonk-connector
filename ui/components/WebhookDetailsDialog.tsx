import React, { useEffect, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { apiClient } from "@/lib/api-client";
import { useToast } from "@/components/ui/use-toast";
import { WebhookLog } from "@/lib/types";

interface WebhookDetailsDialogProps {
  logId: string;
  onClose: () => void;
}

export const WebhookDetailsDialog: React.FC<WebhookDetailsDialogProps> = ({
  logId,
  onClose,
}) => {
  const [log, setLog] = useState<WebhookLog | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [replayLoading, setReplayLoading] = useState(false);
  const { toast } = useToast();

  useEffect(() => {
    const fetchLogDetails = async () => {
      try {
        setLoading(true);
        const response = await apiClient.get(`/webhook-logs/${logId}`);
        setLog(response.data);
      } catch (err) {
        setError("Failed to fetch webhook details");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchLogDetails();
  }, [logId]);

  const handleReplay = async () => {
    if (!log) return;

    try {
      setReplayLoading(true);
      await apiClient.post(`/webhook-logs/${log.id}/replay`);
      toast({
        title: "Webhook Replayed",
        description: "The webhook has been successfully replayed.",
      });
    } catch (err) {
      console.error("Failed to replay webhook:", err);
      toast({
        title: "Replay Failed",
        description: "Failed to replay the webhook. Please try again.",
        variant: "destructive",
      });
    } finally {
      setReplayLoading(false);
    }
  };

  return (
    <Dialog open={true} onOpenChange={onClose}>
      <DialogContent className="max-w-[80vw] max-h-[80vh] flex flex-col p-0">
        <DialogHeader className="px-6 py-4">
          <DialogTitle>Webhook Details</DialogTitle>
        </DialogHeader>
        <ScrollArea className="flex-grow overflow-y-scroll">
          <div className="px-6 py-4">
            {loading ? (
              <WebhookDetailsSkeleton />
            ) : error ? (
              <div className="text-red-500">{error}</div>
            ) : log ? (
              <>
                <Button
                  onClick={handleReplay}
                  disabled={replayLoading}
                  className="mb-4"
                >
                  {replayLoading ? "Replaying..." : "Replay Webhook"}
                </Button>
                <Tabs defaultValue="request" className="w-full">
                  <TabsList>
                    <TabsTrigger value="request">Request</TabsTrigger>
                    <TabsTrigger value="response">Response</TabsTrigger>
                  </TabsList>
                  <TabsContent value="request">
                    <div className="mt-4 space-y-4">
                      <ExpandableSection
                        title="Headers"
                        content={log.headers}
                      />
                      <ExpandableSection title="Body" content={log.body} />
                    </div>
                  </TabsContent>
                  <TabsContent value="response">
                    <div className="mt-4 space-y-4">
                      <h3 className="font-semibold">Status Code</h3>
                      <p>{log.status_code}</p>
                      <ExpandableSection
                        title="Body"
                        content={log.response_body}
                      />
                    </div>
                  </TabsContent>
                </Tabs>
              </>
            ) : null}
          </div>
        </ScrollArea>
      </DialogContent>
    </Dialog>
  );
};

// ... (ExpandableSection and WebhookDetailsSkeleton components remain the same)

interface ExpandableSectionProps {
  title: string;
  content: string;
}

const ExpandableSection: React.FC<ExpandableSectionProps> = ({
  title,
  content,
}) => {
  const [isExpanded, setIsExpanded] = useState(true);

  const toggleExpand = () => setIsExpanded(!isExpanded);

  return (
    <div className="flex flex-col overflow-y-auto">
      <h3 className="font-semibold">{title}</h3>
      <div
        className={`mt-2 bg-muted p-2 rounded max-w-[75vw] overflow-auto ${
          isExpanded ? "" : "max-h-40 overflow-hidden"
        }`}
      >
        <pre className="whitespace-pre-wrap">
          {JSON.stringify(JSON.parse(content), null, 2)}
        </pre>
      </div>
      <Button variant="ghost" onClick={toggleExpand} className="mt-2">
        {isExpanded ? "Show Less" : "Show More"}
      </Button>
    </div>
  );
};

const WebhookDetailsSkeleton: React.FC = () => (
  <div className="space-y-4">
    <Skeleton className="h-4 w-[100px]" />
    <Skeleton className="h-[20px] w-full" />
    <Skeleton className="h-[100px] w-full" />
    <Skeleton className="h-4 w-[100px]" />
    <Skeleton className="h-[100px] w-full" />
  </div>
);
