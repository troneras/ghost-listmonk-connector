import React from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Pencil, Trash2, ToggleLeft, ToggleRight } from "lucide-react";
import { useSonContext } from "@/contexts/SonContext";
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
import { SonListSkeleton } from "./SonListSkeleton";
import { useToast } from "@/components/ui/use-toast";
import { Son } from "@/lib/types";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

function formatDuration(duration: string): string {
  console.log("duration", duration);
  const match = duration.match(/^(\d+)([smhdw])$/);
  if (!match) return duration;

  const [, value, unit] = match;
  const num = parseInt(value, 10);

  switch (unit) {
    case "s":
      return `${num} second${num !== 1 ? "s" : ""}`;
    case "m":
      return `${num} minute${num !== 1 ? "s" : ""}`;
    case "h":
      return `${num} hour${num !== 1 ? "s" : ""}`;
    case "d":
      return `${num} day${num !== 1 ? "s" : ""}`;
    case "w":
      return `${num} week${num !== 1 ? "s" : ""}`;
    default:
      return duration;
  }
}

export function SonList() {
  const { sons, loading, error, deleteSon, updateSon } = useSonContext();
  const { toast } = useToast();

  if (loading) {
    return <SonListSkeleton />;
  }

  if (error) {
    return <div>Error: {error.message}</div>;
  }

  if (sons.length === 0) {
    return <div>No sons found. Create your first son!</div>;
  }

  const handleToggleEnabled = async (son: Son) => {
    try {
      await updateSon(son.id, { ...son, enabled: !son.enabled });
      toast({
        title: "Son Updated",
        description: `${son.name} has been ${
          !son.enabled ? "enabled" : "disabled"
        }.`,
        variant: "default",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to update Son status.",
        variant: "destructive",
      });
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Your Sons</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Trigger</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Delay</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sons.map((son) => (
              <TableRow key={son.id}>
                <TableCell>{son.name}</TableCell>
                <TableCell>{son.trigger}</TableCell>
                <TableCell>
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => handleToggleEnabled(son)}
                        >
                          {son.enabled ? (
                            <ToggleRight className="h-4 w-4 text-green-500" />
                          ) : (
                            <ToggleLeft className="h-4 w-4 text-red-500" />
                          )}
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>
                        <p>{son.enabled ? "Enabled" : "Disabled"}</p>
                      </TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                </TableCell>
                <TableCell>{formatDuration(son.delay)}</TableCell>
                <TableCell>
                  <div className="flex space-x-2">
                    <Link href={`/sons/${son.id}`}>
                      <Button variant="outline" size="icon">
                        <Pencil className="h-4 w-4" />
                      </Button>
                    </Link>
                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button variant="outline" size="icon">
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                          <AlertDialogDescription>
                            This action cannot be undone. This will permanently
                            delete the Son.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Cancel</AlertDialogCancel>
                          <AlertDialogAction onClick={() => deleteSon(son.id)}>
                            Delete
                          </AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
