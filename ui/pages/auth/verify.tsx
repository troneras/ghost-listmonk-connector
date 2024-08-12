import { useEffect } from "react";
import { useRouter } from "next/router";
import { useAuthContext } from "@/contexts/AuthContext";
import { useCustomToast } from "@/hooks/useCustomToast";
import { apiClient } from "@/lib/api-client";

const VerifyPage = () => {
  const router = useRouter();
  const { login } = useAuthContext();
  const { showToast } = useCustomToast();

  useEffect(() => {
    const verifyToken = async () => {
      const { token } = router.query;
      if (token) {
        try {
          const response = await apiClient.get(`/auth/verify?token=${token}`);
          await login(response.data.token, response.data.user);
          showToast("Success", "Logged in successfully");
          router.push("/");
        } catch (error) {
          console.error("Verification failed:", error);
          showToast("Error", "Verification failed", "destructive");
          router.push("/login");
        }
      }
    };

    verifyToken();
  }, [router.query]);

  return (
    <div className="flex items-center justify-center h-screen">
      <p>Verifying your magic link...</p>
    </div>
  );
};

VerifyPage.layout = "minimal";

export default VerifyPage;
