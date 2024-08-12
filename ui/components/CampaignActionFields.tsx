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
import { Textarea } from "@/components/ui/textarea";
import { X } from "lucide-react";
import { CampaignActionFieldsProps } from "@/lib/types";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";

export function CampaignActionFields({
  form,
  index,
  lists,
  templates,
}: CampaignActionFieldsProps) {
  const placeholderTemplate = `
  <h1>New Blog Post: {{ .Post.Title }}</h1>
  
  <img src="{{ .Post.FeatureImage }}" alt="Feature image for {{ .Post.Title }}">
  
  <h2>{{ .Post.CustomExcerpt }}</h2>
  
  <div>
    {{ .Post.Html }}
  </div>
  
  <p>Read the full post at: <a href="https://yourblog.com/{{ .Post.Slug }}">{{ .Post.Title }}</a></p>
  
  <p>Published on: {{ .Post.PublishedAt }}</p>
  `;

  return (
    <div className="space-y-4">
      <FormField
        control={form.control}
        name={`actions.${index}.parameters.name`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Campaign Name</FormLabel>
            <FormControl>
              <Input {...field} placeholder="Enter campaign name" />
            </FormControl>
            <FormDescription>
              Enter a unique name for this campaign.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
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
        name={`actions.${index}.parameters.type`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Campaign Type</FormLabel>
            <Select onValueChange={field.onChange} value={field.value}>
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select campaign type" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectItem value="regular">Regular</SelectItem>
                <SelectItem value="optin">Opt-in</SelectItem>
              </SelectContent>
            </Select>
            <FormDescription>
              Choose the type of campaign you want to create.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name={`actions.${index}.parameters.content_type`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Content Type</FormLabel>
            <Select onValueChange={field.onChange} value={field.value}>
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select content type" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectItem value="richtext">Rich Text</SelectItem>
                <SelectItem value="html">HTML</SelectItem>
                <SelectItem value="markdown">Markdown</SelectItem>
                <SelectItem value="plain">Plain Text</SelectItem>
              </SelectContent>
            </Select>
            <FormDescription>
              Choose the content type for your campaign.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name={`actions.${index}.parameters.body`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Campaign Body</FormLabel>
            <FormControl>
              <Textarea
                {...field}
                placeholder={placeholderTemplate}
                rows={15}
              />
            </FormControl>
            <FormDescription>
              <p>
                The main content of your campaign. Use the format specified in
                the Content Type field.
              </p>
              <p>Available placeholders:</p>
              <ul className="list-disc pl-5 space-y-1">
                <li>
                  <code>{"{{ .Post.Title }}"}</code> - The title of the post
                </li>
                <li>
                  <code>{"{{ .Post.FeatureImage }}"}</code> - URL of the feature
                  image
                </li>
                <li>
                  <code>{"{{ .Post.CustomExcerpt }}"}</code> - The custom
                  excerpt of the post
                </li>
                <li>
                  <code>{"{{ .Post.Html }}"}</code> - The full HTML content of
                  the post
                </li>
                <li>
                  <code>{"{{ .Post.Slug }}"}</code> - The URL slug of the post
                </li>
                <li>
                  <code>{"{{ .Post.PublishedAt }}"}</code> - The publication
                  date of the post
                </li>
              </ul>
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              variant="outline"
              type="button"
              onClick={() =>
                form.setValue(
                  `actions.${index}.parameters.body`,
                  placeholderTemplate
                )
              }
            >
              Insert Template Example
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <p>
              Click to insert a template example using all available
              placeholders
            </p>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>

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
        name={`actions.${index}.parameters.send_at`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Schedule Send</FormLabel>
            <FormControl>
              <DatePicker
                selected={field.value ? new Date(field.value) : null}
                onChange={(date) =>
                  field.onChange(date ? date.toISOString() : null)
                }
                showTimeSelect
                dateFormat="MMMM d, yyyy h:mm aa"
                minDate={new Date(Date.now() + 5 * 60 * 1000)}
                placeholderText="Select date and time"
                className="w-full p-2 border rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                calendarClassName="bg-white border rounded-md shadow-lg"
                wrapperClassName="w-full"
                popperClassName="z-10"
                timeClassName={() => "text-blue-600"}
                dayClassName={() => "hover:bg-blue-100"}
                monthClassName={() => "text-gray-700"}
                weekDayClassName={() => "text-gray-500"}
              />
            </FormControl>
            <FormDescription>
              Schedule when to send this campaign. If not set or set to a time
              in the past, the campaign will be scheduled for 5 minutes from
              now.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />
    </div>
  );
}
