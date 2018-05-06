import { MediaModule } from './media.module';

describe('MediaModule', () => {
  let mediaModule: MediaModule;

  beforeEach(() => {
    mediaModule = new MediaModule();
  });

  it('should create an instance', () => {
    expect(mediaModule).toBeTruthy();
  });
});
