export const prerender = false;

export function load({ params }) {
    return {
        exceptionHash: params.exceptionHash,
        exceptionId: params.exceptionId
    };
}
