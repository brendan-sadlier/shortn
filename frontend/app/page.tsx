import { Button } from '@/components/ui/button';
import { Zap } from 'lucide-react';
import Link from 'next/link';

export default function Home() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-24 gap-y-8">
      <div className="bg-primary text-primary-foreground rounded-2xl p-4 shadow-lg">
        <Zap className="w-10 h-10" />
      </div>

      <h1 className="text-4xl font-bold">Shortn</h1>

      <div className="flex items-center justify-center w-full max-w-2xl gap-x-4">
        <Button asChild variant="outline" size="lg">
          <Link href="/auth/signin">Sign In</Link>
        </Button>

        <Button asChild size="lg">
          <Link href="/auth/signup">Sign Up</Link>
        </Button>
      </div>
    </div>
  );
}
