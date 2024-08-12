import type { NextComponentType, NextPageContext } from 'next'

declare module 'next' {
    export type NextPageWithExtras<P = {}, IP = P> = NextComponentType<NextPageContext, IP, P> & {
        auth?: boolean,
        layout?: 'default' | 'admin' | 'minimal' | (string & {})
        permissions?: string[]
        dataFetching?: 'ssg' | 'ssr' | 'csr'
    }
}

