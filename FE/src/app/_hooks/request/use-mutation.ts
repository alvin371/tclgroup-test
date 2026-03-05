import { useMutation as useTanstackMutation } from '@tanstack/react-query'
import type { UseMutationOptions } from '@tanstack/react-query'
import type { TApiResponseError } from '@/commons/types/api'

export function useMutation<TData, TVariables>(
  options: Omit<UseMutationOptions<TData, TApiResponseError, TVariables>, 'mutationFn'> & {
    mutationFn: (variables: TVariables) => Promise<TData>
  },
) {
  return useTanstackMutation<TData, TApiResponseError, TVariables>(options)
}
