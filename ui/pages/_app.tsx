// ui/pages/_app.tsx
"use client";
import { AppProps } from "next/app";
import type { NextPageWithExtras } from "next";
import localFont from "next/font/local";
import { cn } from "@/lib/utils";
import DefaultLayout from "@/components/layout";
import MinimalLayout from "@/components/MinimalLayout";
import { Toaster } from "@/components/ui/toaster";
import "@/styles/globals.css";

import { ProtectedRoute } from "@/components/ProtectedRoute";
import { AuthProvider } from "@/contexts/AuthContext";
import { AnimatePresence, motion } from "framer-motion";
import { useRouter } from "next/router";
import { useToast } from "@/components/ui/use-toast";
import { useEffect } from "react";

const fontSans = localFont({
  src: "../public/fonts/inter-var-latin.woff2",
  variable: "--font-sans",
});

type AppPropsWithAuth = AppProps & {
  Component: NextPageWithExtras;
};

function MyApp({ Component, pageProps }: AppPropsWithAuth) {
  const router = useRouter();
  const Layout = Component.layout === "minimal" ? MinimalLayout : DefaultLayout;
  const { toast } = useToast();

  useEffect(() => {
    const handleUnauthorized = (event: Event) => {
      if (event instanceof CustomEvent && event.detail === "unauthorized") {
        toast({
          title: "Session Expired",
          description:
            "You've been logged out due to inactivity. Please log in again.",
          variant: "destructive",
        });
      }
    };

    window.addEventListener("unauthorized", handleUnauthorized);

    return () => {
      window.removeEventListener("unauthorized", handleUnauthorized);
    };
  }, [toast]);

  return (
    <div
      className={cn(
        "min-h-screen max-h-screen bg-background font-sans antialiased",
        fontSans.variable
      )}
    >
      <AuthProvider>
        <AnimatePresence mode="wait">
          <Layout>
            <motion.div
              key={router.route}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.2 }}
              className="flex w-full"
            >
              {Component.auth ? (
                <ProtectedRoute>
                  <Component {...pageProps} />
                </ProtectedRoute>
              ) : (
                <Component {...pageProps} />
              )}
            </motion.div>
          </Layout>
        </AnimatePresence>
      </AuthProvider>
      <Toaster />
    </div>
  );
}

export default MyApp;
