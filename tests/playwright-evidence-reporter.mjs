import fs from 'node:fs'
import path from 'node:path'

function sanitizeSegment(value) {
  return String(value)
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

function buildEvidenceName(test, result, appName, attachmentPath) {
  const baseName = [appName, sanitizeSegment(test.title), result.status].filter(Boolean).join('--')
  const extension = path.extname(attachmentPath) || '.bin'
  return `${baseName}${extension}`
}

export default class PlaywrightEvidenceReporter {
  constructor(options = {}) {
    this.appName = options.appName || 'app'
    this.evidenceDir = options.evidenceDir
      ? path.resolve(options.evidenceDir)
      : path.resolve(process.cwd(), '../tests')
  }

  onBegin() {
    fs.mkdirSync(this.evidenceDir, { recursive: true })
  }

  onTestEnd(test, result) {
    for (const attachment of result.attachments) {
      if (!attachment.path) continue
      if (!/\.(webm|png|zip)$/i.test(attachment.path)) continue

      const evidenceName = buildEvidenceName(test, result, this.appName, attachment.path)
      const destinationPath = path.join(this.evidenceDir, evidenceName)
      fs.copyFileSync(attachment.path, destinationPath)
    }
  }
}
