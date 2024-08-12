// Create a new file: ui/components/WebhookInfo.tsx
import React, { useState } from "react";
import { useSonContext } from "@/contexts/SonContext";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Eye, EyeOff } from "lucide-react";

export function WebhookInfo() {
  const { webhook, webhookLoading, webhookError } = useSonContext();
  const [showSecret, setShowSecret] = useState(false);

  if (webhookLoading) {
    return <div>Loading webhook information...</div>;
  }

  if (webhookError) {
    return <div>Error loading webhook information: {webhookError.message}</div>;
  }

  if (!webhook) {
    return <div>No webhook information available.</div>;
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Webhook Information</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <label className="text-sm font-medium">Endpoint:</label>
          <Input value={webhook.endpoint} readOnly />
        </div>
        <div>
          <label className="text-sm font-medium">Secret:</label>
          <div className="flex items-center space-x-2">
            <Input
              type={showSecret ? "text" : "password"}
              value={webhook.secret}
              readOnly
            />
            <Button
              variant="outline"
              size="icon"
              onClick={() => setShowSecret(!showSecret)}
            >
              {showSecret ? (
                <EyeOff className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
