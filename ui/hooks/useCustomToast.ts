import { useToast } from "@/components/ui/use-toast"

export function useCustomToast() {
    const { toast } = useToast()

    const showToast = (title: string, description: string, variant: "default" | "destructive" = "default") => {
        toast({
            title,
            description,
            variant,
        })
    }

    return { showToast }
}