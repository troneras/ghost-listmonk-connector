import React from "react";
import { UseFormReturn } from "react-hook-form";
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
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";
import { CampaignActionFields } from "./CampaignActionFields";
import { ManageSubscriberActionFields } from "./ManageSubscriberActionFields";
import { TransactionalEmailActionFields } from "./TransactionalEmailActionFields";
import { ListmonkList, ListmonkTemplate } from "@/lib/types";

interface ActionFormProps {
  form: UseFormReturn<any>;
  index: number;
  remove: (index: number) => void;
  lists: ListmonkList[];
  templates: ListmonkTemplate[];
}

export function ActionForm({
  form,
  index,
  remove,
  lists,
  templates,
}: ActionFormProps) {
  const actionType: string = form.watch(`actions.${index}.type`);
  const trigger: string = form.watch("trigger");

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
                onValueChange={(value: string) => {
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
              {actionType === "manage_subscriber" && (
                <FormDescription>
                  <p>
                    This action will automatically create the subscriber on
                    listmonk if he doesn't exist already.
                  </p>
                  Subscribers will be automatically added or removed to lists
                  with the same name as the newsletters on Ghost.
                </FormDescription>
              )}

              <FormMessage />
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

        {actionType === "manage_subscriber" && trigger === "member_created" && (
          <ManageSubscriberActionFields
            form={form}
            index={index}
            lists={lists}
          />
        )}

        {actionType === "send_transactional_email" && (
          <TransactionalEmailActionFields
            form={form}
            index={index}
            templates={templates}
          />
        )}
      </CardContent>
    </Card>
  );
}
