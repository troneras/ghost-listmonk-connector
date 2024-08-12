import React, { useState, useEffect } from "react";
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
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Trash2, Plus } from "lucide-react";
import { ListmonkTemplate } from "@/lib/types";
import Link from "next/link";

interface TransactionalEmailActionFieldsProps {
  form: UseFormReturn<any>;
  index: number;
  templates: ListmonkTemplate[];
}

interface Header {
  key: string;
  value: string;
}

interface EditingField {
  key: string;
  newKey: string;
  value: string;
}

export function TransactionalEmailActionFields({
  form,
  index,
  templates,
}: TransactionalEmailActionFieldsProps) {
  const [localData, setLocalData] = useState<Record<string, string>>({});
  const [editingField, setEditingField] = useState<EditingField | null>(null);

  useEffect(() => {
    const formData = form.getValues(`actions.${index}.parameters.data`);
    if (formData && Object.keys(formData).length > 0) {
      setLocalData(formData);
    }
  }, [form, index]);

  const addDataField = () => {
    const newKey = `newField${Object.keys(localData).length}`;
    setLocalData((prev) => ({ ...prev, [newKey]: "" }));
    setEditingField({ key: newKey, newKey: newKey, value: "" });
  };

  const removeDataField = (keyToRemove: string) => {
    setLocalData((prev) => {
      const newData = { ...prev };
      delete newData[keyToRemove];
      return newData;
    });
    setEditingField(null);
  };

  const handleKeyChange = (oldKey: string, newKey: string) => {
    setEditingField((prev) => (prev ? { ...prev, newKey } : null));
  };

  const handleValueChange = (key: string, value: string) => {
    setLocalData((prev) => ({ ...prev, [key]: value }));
  };

  const handleKeyBlur = () => {
    if (editingField && editingField.key !== editingField.newKey) {
      setLocalData((prev) => {
        const newData = { ...prev };
        delete newData[editingField.key];
        newData[editingField.newKey] = editingField.value;
        return newData;
      });
    }
    setEditingField(null);
  };

  useEffect(() => {
    form.setValue(`actions.${index}.parameters.data`, localData);
  }, [localData, form, index]);
  return (
    <>
      {/* Template selection field remains unchanged */}
      <FormField
        control={form.control}
        name={`actions.${index}.parameters.template_id`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Template</FormLabel>
            <Select
              onValueChange={(value: string) => field.onChange(parseInt(value))}
              value={field.value?.toString()}
            >
              <FormControl>
                <SelectTrigger>
                  <SelectValue placeholder="Select a template" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                {templates.map((template) => (
                  <SelectItem key={template.id} value={template.id.toString()}>
                    {template.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormDescription>
              Check the available template expressions in this link -{" "}
              <Link href="https://listmonk.app/docs/templating" target="_blank">
                https://listmonk.app/docs/templating
              </Link>
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name={`actions.${index}.parameters.headers`}
        render={({ field }) => (
          <FormItem>
            <FormLabel>Headers</FormLabel>
            <div className="space-y-2">
              {((field.value as Header[]) || []).map((header, headerIndex) => (
                <div key={headerIndex} className="flex items-center space-x-2">
                  <Input
                    placeholder="Key"
                    value={header.key}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                      const newHeaders = [...(field.value as Header[])];
                      newHeaders[headerIndex].key = e.target.value;
                      field.onChange(newHeaders);
                    }}
                  />
                  <Input
                    placeholder="Value"
                    value={header.value}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                      const newHeaders = [...(field.value as Header[])];
                      newHeaders[headerIndex].value = e.target.value;
                      field.onChange(newHeaders);
                    }}
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => {
                      const newHeaders = (field.value as Header[]).filter(
                        (_, i) => i !== headerIndex
                      );
                      field.onChange(newHeaders);
                    }}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              ))}
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => {
                  const newHeaders = [
                    ...((field.value as Header[]) || []),
                    { key: "", value: "" },
                  ];
                  field.onChange(newHeaders);
                }}
              >
                <Plus className="mr-2 h-4 w-4" /> Add Header
              </Button>
            </div>
            <FormMessage />
            <FormDescription>
              Headers will be added on the email, you can use this to track
              campaigns.
            </FormDescription>
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name={`actions.${index}.parameters.data`}
        render={() => (
          <FormItem>
            <FormLabel>Additional Data</FormLabel>
            <div className="space-y-2">
              {Object.entries(localData).map(([key, value]) => (
                <div key={key} className="flex items-center space-x-2">
                  <Input
                    placeholder="Key"
                    value={
                      editingField?.key === key ? editingField.newKey : key
                    }
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      handleKeyChange(key, e.target.value)
                    }
                    onFocus={() => setEditingField({ key, newKey: key, value })}
                    onBlur={handleKeyBlur}
                  />
                  <Input
                    placeholder="Value"
                    value={value}
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                      handleValueChange(key, e.target.value)
                    }
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => removeDataField(key)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              ))}
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={addDataField}
              >
                <Plus className="mr-2 h-4 w-4" /> Add Data Field
              </Button>
            </div>
            <FormMessage />
            <FormDescription>
              <p>Enter any additional data for the transactional email.</p>
              <p>
                Available in the template as{" "}
                <code className="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
                  &#123;&#123; .Tx.Data.* &#125;&#125;
                </code>
              </p>
            </FormDescription>
          </FormItem>
        )}
      />
    </>
  );
}
