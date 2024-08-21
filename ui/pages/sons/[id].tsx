import { NextPageWithExtras } from "next";

import React, { useEffect } from "react";
import { useRouter } from "next/router";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Form } from "@/components/ui/form";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PlusCircle } from "lucide-react";
import { useSon } from "@/hooks/useSon";
import { useCustomToast } from "@/hooks/useCustomToast";
import { SonDetailsForm } from "@/components/SonDetailsForm";
import { ActionForm } from "@/components/ActionForm";
import { editableSonSchema } from "@/lib/schemas";
import { EditableSon } from "@/lib/types";
import { useSonContext } from "@/contexts/SonContext";
import { useListContext } from "@/contexts/ListContext";
import { useTemplateContext } from "@/contexts/TemplateContext";

function isValidTrigger(trigger: string): trigger is EditableSon["trigger"] {
  return [
    "member_created",
    "member_deleted",
    "member_updated",
    "post_published",
    "post_scheduled",
  ].includes(trigger);
}

const SonDetailPage: NextPageWithExtras = () => {
  const router = useRouter();
  const { id } = router.query;
  const { son, loading: sonLoading, error: sonError } = useSon(id as string);
  const { updateSon, deleteSon } = useSonContext();
  const { showToast } = useCustomToast();
  const { lists, loading: listsLoading, error: listsError } = useListContext();
  const {
    templates,
    loading: templatesLoading,
    error: templatesError,
  } = useTemplateContext();

  const form = useForm<EditableSon>({
    resolver: zodResolver(editableSonSchema),
    defaultValues: {
      name: "",
      trigger: "member_created",
      delay: "0",
      actions: [],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "actions",
  });

  useEffect(() => {
    if (son) {
      const { name, trigger, delay, actions } = son;

      if (isValidTrigger(trigger)) {
        form.reset({
          name,
          trigger,
          delay,
          actions: actions.map((action) => ({
            type: action.type,
            parameters: action.parameters,
          })),
        });
      } else {
        console.error(`Invalid trigger value: ${trigger}`);
        form.reset({
          name,
          trigger: "member_created",
          delay,
          actions: actions.map((action) => ({
            type: action.type,
            parameters: action.parameters,
          })),
        });
      }
    }
  }, [son, form]);

  if (sonLoading || listsLoading || templatesLoading) {
    return (
      <div className="flex justify-center items-center h-full">Loading...</div>
    );
  }

  if (sonError || listsError || templatesError) {
    return (
      <div className="text-red-500">
        Error:{" "}
        {sonError?.message || listsError?.message || templatesError?.message}
      </div>
    );
  }

  if (!son) return <div>Son not found</div>;

  const onSubmit = async (data: EditableSon) => {
    try {
      await updateSon(son.id, data);
      showToast("Success", "Son updated successfully");
      router.push("/sons");
    } catch (error) {
      showToast("Error", "Failed to update Son", "destructive");
      console.error("Failed to update Son:", error);
    }
  };

  const handleDelete = async () => {
    try {
      await deleteSon(son.id);
      showToast("Success", "Son deleted successfully");
      router.push("/sons");
    } catch (error) {
      showToast("Error", "Failed to delete Son", "destructive");
      console.error("Failed to delete Son:", error);
    }
  };

  return (
    <div className="flex-1 max-w-3xl mx-auto space-y-4">
      <h1 className="text-2xl font-bold">Edit Son: {son.name}</h1>
      <p className="text-sm text-gray-500">
        Created at: {new Date(son.created_at).toLocaleString()}
      </p>
      <p className="text-sm text-gray-500">
        Last updated: {new Date(son.updated_at).toLocaleString()}
      </p>

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
            Save Changes
          </Button>
        </form>
      </Form>
      <AlertDialog>
        <AlertDialogTrigger asChild>
          <Button variant="destructive">Delete Son</Button>
        </AlertDialogTrigger>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the
              Son.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleDelete}>Delete</AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
};

SonDetailPage.auth = true;

export default SonDetailPage;
