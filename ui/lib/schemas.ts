import * as z from 'zod';

// Define the schema for action parameters
const actionParametersSchema = z.object({
    subject: z.string().optional(),
    lists: z.array(z.number()).optional(),
    template_id: z.number().optional(),
    tags: z.string().optional(),
    send_now: z.boolean().optional(),
}).strict().or(z.record(z.any())); // Allow any other properties for flexibility



// Define the schema for a single action
const actionSchema = z.object({
    type: z.enum(['send_transactional_email', 'manage_subscriber', 'create_campaign']),
    parameters: actionParametersSchema,
});

// Define the schema for the entire Son object
export const sonSchema = z.object({
    id: z.string().optional(), // Optional because it might not be present when creating a new Son
    name: z.string().min(1, 'Name is required'),
    trigger: z.enum([
        'member_created',
        'member_deleted',
        'member_updated',
        'post_published',
        'post_scheduled',
    ]),
    delay: z.number().min(0, 'Delay must be a positive number'),
    actions: z.array(actionSchema).min(1, 'At least one action is required'),
    created_at: z.date().optional(),
    updated_at: z.date().optional(),
});

export const editableSonSchema = z.object({
    name: z.string().min(1, 'Name is required'),
    trigger: z.enum([
        'member_created',
        'member_deleted',
        'member_updated',
        'post_published',
        'post_scheduled',
    ]),
    delay: z.number().min(0, 'Delay must be a positive number'),
    actions: z.array(z.object({
        type: z.enum(['send_transactional_email', 'manage_subscriber', 'create_campaign']),
        parameters: z.record(z.any()),
    })),
});


// Define a schema for creating a new Son (without id, createdAt, and updatedAt)
export const createSonSchema = sonSchema.omit({ id: true, created_at: true, updated_at: true });

// Define a schema for updating an existing Son
export const updateSonSchema = sonSchema.partial().extend({
    id: z.string(),
});

// Define types for create and update operations
export type CreateSonInput = z.infer<typeof createSonSchema>;
export type UpdateSonInput = z.infer<typeof updateSonSchema>;