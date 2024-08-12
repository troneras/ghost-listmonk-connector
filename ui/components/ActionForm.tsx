import React from "react";
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";
import { CampaignActionFields } from "./CampaignActionFields";
import { ManageSubscriberActionFields } from "./ManageSubscriberActionFields";
import { ActionFormProps } from "@/lib/types";

export function ActionForm({
  form,
  index,
  remove,
  lists,
  templates,
}: ActionFormProps) {
  const actionType = form.watch(`actions.${index}.type`);
  const trigger = form.watch("trigger");

  return (
    <Card className="border border-gray-200">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">
          Action {index + 1}
        </CardTitle>
        <Button
          type="button"
          onClick={() => remove(index)}
          variant="ghost"
          size="sm"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </CardHeader>
      <CardContent className="space-y-4">
        <FormField
          control={form.control}
          name={`actions.${index}.type`}
          render={({ field }) => (
            <FormItem>
              <FormLabel>Action Type</FormLabel>
              <Select
                onValueChange={(value) => {
                  field.onChange(value);
                  form.setValue(`actions.${index}.parameters`, {});
                }}
                value={field.value}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select an action type" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="send_transactional_email">
                    Send Transactional Email
                  </SelectItem>
                  <SelectItem value="manage_subscriber">
                    Manage Subscriber
                  </SelectItem>
                  <SelectItem value="create_campaign">
                    Create Campaign
                  </SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
              {actionType === "manage_subscriber" && (
                <FormDescription>
                  Subscriber management automatically adds or removes
                  subscribers on listmonk and adds to lists based on the trigger
                  event and newsletters from ghost blog.
                </FormDescription>
              )}
            </FormItem>
          )}
        />

        {actionType === "create_campaign" && (
          <CampaignActionFields
            form={form}
            index={index}
            lists={lists}
            templates={templates}
          />
        )}

        {actionType === "manage_subscriber" && trigger == "member_created" && (
          <ManageSubscriberActionFields
            form={form}
            index={index}
            lists={lists}
          />
        )}

        {actionType === "send_transactional_email" && (
          <FormField
            control={form.control}
            name={`actions.${index}.parameters`}
            render={({ field }) => (
              <FormItem>
                <FormLabel>Parameters</FormLabel>
                <FormControl>
                  <Textarea
                    {...field}
                    onChange={(e) => {
                      try {
                        const parsedValue = JSON.parse(e.target.value);
                        field.onChange(parsedValue);
                      } catch (error) {
                        field.onChange(e.target.value);
                      }
                    }}
                    value={
                      typeof field.value === "object"
                        ? JSON.stringify(field.value, null, 2)
                        : field.value
                    }
                    className="font-mono text-sm"
                    rows={5}
                  />
                </FormControl>
                <FormDescription>
                  Enter the parameters for the transactional email as JSON.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        )}
      </CardContent>
    </Card>
  );
}
