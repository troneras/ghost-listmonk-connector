import React from "react";
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { X } from "lucide-react";
import { CampaignActionFieldsProps } from "@/lib/types";

export function CampaignActionFields({ form, index, lists, templates }: CampaignActionFieldsProps) {
  return (
    <div className="space-y-4">
      <FormField
        control={form.control}
        name={`actions.${index}.parameters.subject`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Subject</FormLabel>
            <FormControl>
              <Input {...field} placeholder="Enter campaign subject" />
            </FormControl>
            <FormDescription>
              The subject line for your campaign email.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name={`actions.${index}.parameters.lists`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Lists</FormLabel>
            <FormControl>
              <Select
                onValueChange={(value) => {
                  const newLists = [...(field.value || []), parseInt(value)];
                  field.onChange(newLists);
                }}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select lists" />
                </SelectTrigger>
                <SelectContent>
                  {lists.map((list) => (
                    <SelectItem key={list.id} value={list.id.toString()}>
                      {list.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </FormControl>
            <div className="mt-2 flex flex-wrap gap-2">
              {field.value?.map((listId: number) => {
                const list = lists.find((l) => l.id === listId);
                return (
                  <Badge key={listId} variant="secondary" className="px-2 py-1">
                    {list ? list.name : `List ${listId}`}
                    <Button
                      variant="ghost"
                      size="sm"
                      className="ml-1 h-4 w-4 p-0"
                      onClick={() => {
                        const newLists = field.value.filter(
                          (id: number) => id !== listId
                        );
                        field.onChange(newLists);
                      }}
                    >
                      <X className="h-3 w-3" />
                    </Button>
                  </Badge>
                );
              })}
            </div>
            <FormDescription>
              Select the lists to send this campaign to.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name={`actions.${index}.parameters.template_id`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Template</FormLabel>
            <FormControl>
              <Select
                onValueChange={(value) => field.onChange(parseInt(value))}
                value={field.value?.toString()}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select template" />
                </SelectTrigger>
                <SelectContent>
                  {templates.map((template) => (
                    <SelectItem
                      key={template.id}
                      value={template.id.toString()}
                    >
                      {template.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </FormControl>
            <FormDescription>
              Choose the template for your campaign email.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name={`actions.${index}.parameters.tags`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Tags</FormLabel>
            <FormControl>
              <Input {...field} placeholder="Enter comma-separated tags" />
            </FormControl>
            <FormDescription>
              Add tags to categorize your campaign (comma-separated).
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      <FormField
        control={form.control}
        name={`actions.${index}.parameters.send_now`}
        render={({ field }) => (
          <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
            <div className="space-y-0.5">
              <FormLabel className="text-base">Send Now</FormLabel>
              <FormDescription>
                Toggle to send the campaign immediately upon creation.
              </FormDescription>
            </div>
            <FormControl>
              <Switch checked={field.value} onCheckedChange={field.onChange} />
            </FormControl>
          </FormItem>
        )}
      />
    </div>
  );
}
