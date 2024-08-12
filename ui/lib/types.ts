import { UseFormReturn, FieldValues } from 'react-hook-form';
import { editableSonSchema } from '@/lib/schemas';  // Assuming this is where your Son type is defined
import { z } from 'zod';
// Define a type for the Son object based on the schema
export interface Son extends z.infer<typeof editableSonSchema> {
    id: string;
    created_at: string;
    updated_at: string;
}

// Type for editable Son fields
export type EditableSon = z.infer<typeof editableSonSchema>;

export interface ListmonkList {
    id: number;
    name: string;
    type: string;
    optin: string;
    tags: string[];
}

export interface ListmonkTemplate {
    id: number;
    name: string;
    type: string;
    is_default: boolean;
}
// Props for CampaignActionFields
export interface CampaignActionFieldsProps {
    form: UseFormReturn<EditableSon>;
    index: number;
    lists: ListmonkList[];
    templates: ListmonkTemplate[];
}

// Props for ManageSubscriberActionFields
export interface ManageSubscriberActionFieldsProps {
    form: UseFormReturn<EditableSon>;
    index: number;
    lists: ListmonkList[];
}

// Props for SonDetailsForm
export interface SonDetailsFormProps {
    form: UseFormReturn<EditableSon>;
}

export interface ActionFormProps {
    form: UseFormReturn<EditableSon>;
    index: number;
    remove: (index: number) => void;
    lists: ListmonkList[];
    templates: ListmonkTemplate[];
}