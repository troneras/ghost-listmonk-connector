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
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { X } from "lucide-react";
import { ManageSubscriberActionFieldsProps } from "@/lib/types";

export function ManageSubscriberActionFields({
  form,
  index,
  lists,
}: ManageSubscriberActionFieldsProps) {
  return (
    <div className="space-y-4">
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
              Default lists to add this user on creation. (Apart from the lists
              for the ghost blog newsletters) p.e. "New Subscribers"
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
      {/* Add more fields specific to manage_subscriber action if needed */}
    </div>
  );
}
