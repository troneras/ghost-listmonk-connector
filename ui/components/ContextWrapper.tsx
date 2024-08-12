import { SonProvider } from "@/contexts/SonContext";
import { ListProvider } from "@/contexts/ListContext";
import { TemplateProvider } from "@/contexts/TemplateContext";

export default function ContextWrapper({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SonProvider>
      <ListProvider>
        <TemplateProvider>{children}</TemplateProvider>
      </ListProvider>
    </SonProvider>
  );
}
