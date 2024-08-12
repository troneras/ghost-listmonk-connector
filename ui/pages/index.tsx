import { NextPageWithExtras } from "next";
import Dashboard from "@/components/Dashboard";

const HomePage: NextPageWithExtras = () => {
  return <Dashboard />;
};

HomePage.auth = true;

export default HomePage;
