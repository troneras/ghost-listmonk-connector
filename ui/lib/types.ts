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

export interface Webhook {
    id: string;
    endpoint: string;
    secret: string;
}

export interface WebhookLog {
    id: string;
    user_id: string;
    method: string;
    path: string;
    headers: string;
    body: string;
    status_code: number;
    response_body: string;
    timestamp: string;
    duration: number;
}

export interface SonExecutionLog {
    id: string;
    son_id: string;
    webhook_log_id: string;
    status: 'pending' | 'success' | 'failure' | 'warning';
    executed_at: string;
    error_message: string | null;
}

export interface ActionExecutionLog {
    id: string;
    son_execution_log_id: string;
    action_type: string;
    status: string;
    executed_at: string;
    error_message: string | null;
}

export interface Pagination {
    total: number;
    limit: number;
    offset: number;
    next_offset: number;
}

export interface RecentActivity {
    id: string;
    user_id: string;
    action_type: string;
    description: string;
    timestamp: string;
}

export interface SonStats {
    name: string;
    executions: number;
    success: number;
    failure: number;
}