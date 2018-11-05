export function getAnonymousCredentials(): Credentials {
  return new Credentials('anon', ['ANONYMOUS']);
}

export class Credentials {
  constructor(
    public readonly username: string,
    public readonly roles: string[]
  ) {}

  hasRole(role: string): boolean {
    return this.roles.includes(role);
  }

  isAnonymous() {
    return this.hasRole('ANONYMOUS');
  }
}
