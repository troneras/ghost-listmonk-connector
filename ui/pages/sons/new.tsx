import { NextPageWithExtras } from "next";
import SonCreationForm from "@/components/SonCreationForm";

const NewSonPage: NextPageWithExtras = () => {
  return (
    <div className="flex-1 max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">Create New Son</h1>
      <SonCreationForm />
    </div>
  );
};

NewSonPage.auth = true;

export default NewSonPage;
