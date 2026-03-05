import { useNavigate, useLocation } from 'react-router'
import type { TRoutePath } from '@/commons/route'

export function useRouter() {
  const navigate = useNavigate()
  const location = useLocation()

  return {
    push: (path: TRoutePath) => navigate(path),
    replace: (path: TRoutePath) => navigate(path, { replace: true }),
    back: () => navigate(-1),
    pathname: location.pathname,
  }
}
