import * as z from 'zod';

// Define the schema for action parameters
const actionParametersSchema = z.object({
    subject: z.string().optional(),
    lists: z.array(z.number()).optional(),
    template_id: z.number().optional(),
    tags: z.string().optional(),
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
    delay: z.string().refine((val) => {
        const durationRegex = /^(\d+)\s*(s|m|h|d|w)$/;
        return durationRegex.test(val);
    }, {
        message: "Invalid duration format. Use format like '30m', '2h', '1d', or '1w'.",
    }).default('0s'),
    actions: z.array(actionSchema).min(1, 'At least one action is required'),
    enabled: z.boolean().default(true),
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
    delay: z.string().min(2).default('0s').refine((val) => {
        const durationRegex = /^(\d+)\s*(s|m|h|d|w)$/;
        return durationRegex.test(val);
    }, {
        message: "Invalid duration format. Use format like '30m', '2h', '1d', or '1w'.",
    }),
    actions: z.array(z.object({
        type: z.enum(['send_transactional_email', 'manage_subscriber', 'create_campaign']),
        parameters: z.record(z.any()),
    })),
    enabled: z.boolean().default(true),
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