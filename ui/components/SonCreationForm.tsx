"use client";
import React from "react";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Form } from "@/components/ui/form";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PlusCircle } from "lucide-react";
import { useSons } from "@/hooks/useSons";
import { useLists } from "@/hooks/useLists";
import { useTemplates } from "@/hooks/useTemplates";
import { useCustomToast } from "@/hooks/useCustomToast";
import { SonDetailsForm } from "./SonDetailsForm";
import { ActionForm } from "./ActionForm";

// Import or define your schema here
import { sonSchema } from "@/lib/schemas";

type SonFormValues = z.infer<typeof sonSchema>;

export default function SonCreationForm() {
  const { createSon } = useSons();
  const { showToast } = useCustomToast();
  const router = useRouter();
  const { lists, loading: listsLoading, error: listsError } = useLists();
  const {
    templates,
    loading: templatesLoading,
    error: templatesError,
  } = useTemplates();

  const form = useForm<SonFormValues>({
    resolver: zodResolver(sonSchema),
    defaultValues: {
      name: "",
      trigger: "member_created",
      delay: 0,
      actions: [
        {
          type: "send_transactional_email",
          parameters: {},
        },
      ],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "actions",
  });

  const onSubmit = async (data: SonFormValues) => {
    try {
      await createSon(data);
      showToast("Success", "Son created successfully");
      router.push("/sons");
    } catch (error) {
      showToast("Error", "Failed to create Son", "destructive");
      console.error("Failed to create Son:", error);
    }
  };

  if (listsLoading || templatesLoading) {
    return (
      <div className="flex justify-center items-center h-full">Loading...</div>
    );
  }

  if (listsError || templatesError) {
    return (
      <div className="text-red-500">
        Error: {listsError?.message || templatesError?.message}
      </div>
    );
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <SonDetailsForm form={form} />

        <Card>
          <CardHeader>
            <CardTitle>Actions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            {fields.map((field, index) => (
              <ActionForm
                key={field.id}
                form={form}
                index={index}
                remove={remove}
                lists={lists}
                templates={templates}
              />
            ))}

            <Button
              type="button"
              onClick={() =>
                append({ type: "send_transactional_email", parameters: {} })
              }
              variant="outline"
              className="w-full"
            >
              <PlusCircle className="mr-2 h-4 w-4" /> Add Action
            </Button>
          </CardContent>
        </Card>

        <Button type="submit" className="w-full">
          Create Son
        </Button>
      </form>
    </Form>
  );
}
