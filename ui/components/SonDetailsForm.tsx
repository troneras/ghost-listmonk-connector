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
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { SonDetailsFormProps } from "@/lib/types";

export function SonDetailsForm({ form }: SonDetailsFormProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Son Details</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Son Name</FormLabel>
              <FormControl>
                <Input placeholder="Enter Son name" {...field} />
              </FormControl>
              <FormDescription>
                Give your Son a unique and descriptive name.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="trigger"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Trigger</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select a trigger" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="member_created">Member Created</SelectItem>
                  <SelectItem value="member_deleted">Member Deleted</SelectItem>
                  <SelectItem value="member_updated">Member Updated</SelectItem>
                  <SelectItem value="post_published">Post Published</SelectItem>
                  <SelectItem value="post_scheduled">Post Scheduled</SelectItem>
                </SelectContent>
              </Select>
              <FormDescription>
                Choose the event that will trigger this Son.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="delay"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Delay (minutes)</FormLabel>
              <FormControl>
                <Input
                  type="number"
                  {...field}
                  onChange={(e) => field.onChange(parseInt(e.target.value))}
                />
              </FormControl>
              <FormDescription>
                Set a delay before the Son executes its actions.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      </CardContent>
    </Card>
  );
}
