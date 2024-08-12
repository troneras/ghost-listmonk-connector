import { useState } from "react";
import { useRouter } from "next/router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { useCustomToast } from "@/hooks/useCustomToast";
import { apiClient } from "@/lib/api-client";

const LoginPage = () => {
  const [email, setEmail] = useState("");
  const router = useRouter();
  const { showToast } = useCustomToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await apiClient.post("/auth/magic-link", { email });
      showToast("Success", "Magic link sent to your email");
    } catch (error) {
      console.error("Failed to send magic link:", error);
      showToast("Error", "Failed to send magic link", "destructive");
    }
  };

  return (
    <div className="flex flex-col w-full items-center justify-center bg-gray-100">
      <h1>
        <span className="font-bold text-2xl">Ghost-Listmonk Connector</span>
      </h1>
      <Card className="w-full max-w-md mt-12">
        <CardHeader>
          <CardTitle>Login</CardTitle>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            <Input
              type="email"
              placeholder="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </CardContent>
          <CardFooter>
            <Button type="submit" className="w-full">
              Send Magic Link
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
};

LoginPage.layout = "minimal";

export default LoginPage;
