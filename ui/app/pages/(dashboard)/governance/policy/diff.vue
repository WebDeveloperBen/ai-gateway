<script setup lang="ts">
import type * as Monaco from "monaco-editor"

interface DiffSegment {
  type: "unchanged" | "added" | "removed" | "modified"
  baseText?: string
  comparisonText?: string
  context: string // e.g., "Section 2.1: Authentication Requirements"
}

interface PolicyPair {
  id: string
  type: "aligned" | "unmatched_base" | "unmatched_comparison"
  title: string
  section: string
  baseContent?: string
  comparisonContent?: string
  diffSegments?: DiffSegment[]
  filename: string
  baseVersion?: string
  comparisonVersion?: string
  sourceTrace: string
  similarities: string[]
  keyDifferences: string[]
}

interface PolicyDiffData {
  metadata: {
    runId: string
    timestamp: string
    baseFilename: string
    comparisonFilename: string
    baseVersion: string
    comparisonVersion: string
    summary: string
  }
  pairs: PolicyPair[]
}

const selectedPairId = ref<string | null>(null)
const isLoading = ref(true)
const diffData = ref<PolicyDiffData | null>(null)
const viewMode = ref<"sections" | "overview">("overview")

const mockData: PolicyDiffData = {
  metadata: {
    runId: "diff_2025_03_15_14_30_45_abc123",
    timestamp: "2025-03-15T14:30:45Z",
    baseFilename: "enterprise-data-governance-policy-v1.0.pdf",
    comparisonFilename: "enterprise-data-governance-policy-v2.0.pdf",
    baseVersion: "1.0.0",
    comparisonVersion: "2.0.0",
    summary:
      "Comprehensive enterprise data governance policy revision addressing GDPR compliance, AI governance, and enhanced security controls across 23 policy sections"
  },
  pairs: [
    // PAGE 1 SECTIONS
    {
      id: "pair-1",
      type: "aligned",
      title: "Purpose and Scope",
      section: "Section 1.1: Policy Foundation",
      baseContent:
        "This **Data Governance Policy** establishes the framework for managing, protecting, and utilizing data assets across the organization. The policy applies to all employees, contractors, and third-party vendors who handle organizational data.\n\nThis policy covers all data types including *personal data*, *financial records*, *intellectual property*, and *operational data* regardless of format or storage location.",
      comparisonContent:
        "This **Data Governance Policy** establishes the comprehensive framework for managing, protecting, and utilizing data assets across the organization and its subsidiaries. The policy applies to all employees, contractors, third-party vendors, and `automated systems` who handle or process organizational data.\n\nThis policy covers all data types including *personal data*, *financial records*, *intellectual property*, *operational data*, and `AI training datasets` regardless of format, storage location, or processing method.",
      diffSegments: [
        {
          type: "modified",
          baseText: "establishes the framework for managing",
          comparisonText: "establishes the comprehensive framework for managing",
          context: "Policy scope expansion"
        },
        {
          type: "modified",
          baseText: "across the organization",
          comparisonText: "across the organization and its subsidiaries",
          context: "Organizational coverage"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "and automated systems",
          context: "Coverage extension"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "and AI training datasets",
          context: "Data type addition"
        }
      ],
      filename: "enterprise-data-governance-policy.pdf",
      baseVersion: "1.0",
      comparisonVersion: "2.0",
      sourceTrace: "Section 1.1, Page 1",
      similarities: [
        "Maintains core policy framework structure",
        "Preserves fundamental coverage of all employee types",
        "Continues comprehensive data type coverage"
      ],
      keyDifferences: [
        "Extended scope to include subsidiaries",
        "Added coverage for automated systems",
        "Included AI training datasets in scope"
      ]
    },
    {
      id: "pair-2",
      type: "aligned",
      title: "Data Classification Standards",
      section: "Section 1.2: Information Taxonomy",
      baseContent:
        "All organizational data must be classified according to sensitivity levels:\n\n**Public**: Information that can be freely shared without risk\n**Internal**: Information for internal use only\n**Confidential**: Sensitive information requiring restricted access\n**Restricted**: Highly sensitive information requiring executive approval for access",
      comparisonContent:
        "All organizational data must be classified according to enhanced sensitivity levels:\n\n**Public**: Information that can be freely shared without organizational risk\n**Internal**: Information for internal use only with usage tracking\n**Confidential**: Sensitive information requiring restricted access and audit trails\n**Restricted**: Highly sensitive information requiring executive approval and multi-factor authentication\n**Critical**: `Mission-critical data` requiring board-level oversight and encryption at rest and in transit",
      diffSegments: [
        {
          type: "modified",
          baseText: "according to sensitivity levels",
          comparisonText: "according to enhanced sensitivity levels",
          context: "Classification enhancement"
        },
        {
          type: "modified",
          baseText: "without risk",
          comparisonText: "without organizational risk",
          context: "Risk specification"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "with usage tracking",
          context: "Internal data monitoring"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "and audit trails",
          context: "Confidential data tracking"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "and multi-factor authentication",
          context: "Restricted access enhancement"
        },
        {
          type: "added",
          baseText: "",
          comparisonText:
            "Critical: Mission-critical data requiring board-level oversight and encryption at rest and in transit",
          context: "New classification level"
        }
      ],
      filename: "enterprise-data-governance-policy.pdf",
      baseVersion: "1.0",
      comparisonVersion: "2.0",
      sourceTrace: "Section 1.2, Page 1",
      similarities: [
        "Maintains four-tier classification system foundation",
        "Preserves public information sharing principles",
        "Continues executive approval for restricted data"
      ],
      keyDifferences: [
        "Added fifth 'Critical' classification tier",
        "Enhanced security requirements for each level",
        "Introduced mandatory encryption for critical data",
        "Added usage tracking and audit trail requirements"
      ]
    },
    {
      id: "pair-3",
      type: "aligned",
      title: "Data Subject Rights",
      section: "Section 2.1: Individual Privacy Rights",
      baseContent:
        "Individuals have the following rights regarding their personal data:\n\n• **Right to Access**: Request copies of personal data held\n• **Right to Rectification**: Correct inaccurate personal data\n• **Right to Erasure**: Request deletion of personal data\n• **Right to Restrict Processing**: Limit how personal data is used\n• **Right to Data Portability**: Receive personal data in machine-readable format\n\nAll rights requests must be acknowledged within *48 hours* and fulfilled within *30 days*.",
      comparisonContent:
        "Individuals have the following comprehensive rights regarding their personal data:\n\n• **Right to Access**: Request copies of personal data held with full processing history\n• **Right to Rectification**: Correct inaccurate personal data with verification\n• **Right to Erasure**: Request deletion of personal data unless legal retention required\n• **Right to Restrict Processing**: Limit how personal data is used with granular controls\n• **Right to Data Portability**: Receive personal data in standardized machine-readable format\n• **Right to Object**: `Object to automated decision-making and profiling`\n• **Right to Human Review**: `Request human oversight of automated decisions`\n\nAll rights requests must be acknowledged within *24 hours* and fulfilled within *21 days* with detailed response documentation.",
      diffSegments: [
        {
          type: "modified",
          baseText: "following rights",
          comparisonText: "following comprehensive rights",
          context: "Rights scope enhancement"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "with full processing history",
          context: "Access right enhancement"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "with verification",
          context: "Rectification enhancement"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "unless legal retention required",
          context: "Erasure limitation"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "with granular controls",
          context: "Processing restriction enhancement"
        },
        {
          type: "modified",
          baseText: "machine-readable format",
          comparisonText: "standardized machine-readable format",
          context: "Portability standardization"
        },
        {
          type: "added",
          baseText: "",
          comparisonText:
            "Right to Object: Object to automated decision-making and profiling\nRight to Human Review: Request human oversight of automated decisions",
          context: "New AI-related rights"
        },
        {
          type: "modified",
          baseText: "within 48 hours and fulfilled within 30 days",
          comparisonText: "within 24 hours and fulfilled within 21 days with detailed response documentation",
          context: "Response time improvement"
        }
      ],
      filename: "enterprise-data-governance-policy.pdf",
      baseVersion: "1.0",
      comparisonVersion: "2.0",
      sourceTrace: "Section 2.1, Page 1",
      similarities: [
        "Maintains core GDPR rights framework",
        "Preserves fundamental data subject protections",
        "Continues structured request handling process"
      ],
      keyDifferences: [
        "Added AI-specific rights (objection, human review)",
        "Faster response times (24h acknowledgment, 21 days fulfillment)",
        "Enhanced verification and documentation requirements",
        "More granular control options for data subjects"
      ]
    },
    {
      id: "pair-4",
      type: "aligned",
      title: "Data Retention and Disposal",
      section: "Section 2.2: Data Lifecycle Management",
      baseContent:
        "Personal data shall be retained for **no longer than necessary** for the specified purposes. Standard retention periods:\n\n• *Customer data*: 5 years after last interaction\n• *Employee records*: 7 years after termination\n• *Financial records*: 7 years from fiscal year end\n• *Marketing data*: 2 years after opt-out\n\nData disposal must be conducted using *secure deletion methods* and documented in disposal logs.",
      comparisonContent:
        "Personal data shall be retained for **no longer than necessary** for the specified purposes with `automated lifecycle management`. Standard retention periods:\n\n• *Customer data*: 3 years after last interaction with `automatic purging`\n• *Employee records*: 7 years after termination with `graduated access restrictions`\n• *Financial records*: 7 years from fiscal year end (regulatory requirement)\n• *Marketing data*: 18 months after opt-out with `quarterly reviews`\n• *`AI training data`*: `2 years maximum with model lineage tracking`\n\nData disposal must be conducted using *cryptographic erasure methods* and documented in `immutable disposal logs` with `third-party verification`.",
      diffSegments: [
        {
          type: "added",
          baseText: "",
          comparisonText: "with automated lifecycle management",
          context: "Automation addition"
        },
        {
          type: "modified",
          baseText: "5 years after last interaction",
          comparisonText: "3 years after last interaction with automatic purging",
          context: "Customer data retention reduction"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "with graduated access restrictions",
          context: "Employee record access controls"
        },
        {
          type: "modified",
          baseText: "2 years after opt-out",
          comparisonText: "18 months after opt-out with quarterly reviews",
          context: "Marketing data retention adjustment"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "AI training data: 2 years maximum with model lineage tracking",
          context: "New AI data category"
        },
        {
          type: "modified",
          baseText: "secure deletion methods",
          comparisonText: "cryptographic erasure methods",
          context: "Enhanced disposal security"
        },
        {
          type: "modified",
          baseText: "documented in disposal logs",
          comparisonText: "documented in immutable disposal logs with third-party verification",
          context: "Enhanced disposal documentation"
        }
      ],
      filename: "enterprise-data-governance-policy.pdf",
      baseVersion: "1.0",
      comparisonVersion: "2.0",
      sourceTrace: "Section 2.2, Page 1-2",
      similarities: [
        "Maintains data minimization principle",
        "Preserves regulatory compliance for financial records",
        "Continues secure disposal requirements"
      ],
      keyDifferences: [
        "Reduced customer data retention from 5 to 3 years",
        "Added automated lifecycle management system",
        "Introduced AI training data category with specific rules",
        "Enhanced disposal security with cryptographic methods",
        "Added third-party verification for disposal processes"
      ]
    },
    // PAGE 2 SECTIONS
    {
      id: "pair-5",
      type: "aligned",
      title: "Cross-Border Data Transfers",
      section: "Section 3.1: International Data Flows",
      baseContent:
        "Personal data transfers outside the organization's home jurisdiction require:\n\n• **Adequacy Decision**: Transfer to countries with adequate data protection\n• **Standard Contractual Clauses**: Use of approved legal frameworks\n• **Binding Corporate Rules**: For intra-group transfers\n• **Explicit Consent**: When other safeguards are not available\n\nAll transfers must be logged and include *data mapping documentation*.",
      comparisonContent:
        "Personal data transfers outside the organization's home jurisdiction require enhanced protection measures:\n\n• **Adequacy Decision**: Transfer to countries with adequate data protection as per `current regulatory assessments`\n• **Standard Contractual Clauses**: Use of `latest approved legal frameworks` with `supplementary measures`\n• **Binding Corporate Rules**: For intra-group transfers with `annual compliance reviews`\n• **Explicit Consent**: When other safeguards are not available with `granular purpose specification`\n• **`Data Localization Compliance`**: `Adherence to local data residency requirements`\n• **`Transfer Impact Assessment`**: `Mandatory evaluation of destination country privacy laws`\n\nAll transfers must be logged in `real-time monitoring systems` and include *comprehensive data mapping documentation* with `risk assessments`.",
      diffSegments: [
        {
          type: "added",
          baseText: "",
          comparisonText: "enhanced protection measures",
          context: "Enhanced security approach"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "as per current regulatory assessments",
          context: "Dynamic adequacy evaluation"
        },
        {
          type: "modified",
          baseText: "approved legal frameworks",
          comparisonText: "latest approved legal frameworks with supplementary measures",
          context: "Enhanced SCC requirements"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "with annual compliance reviews",
          context: "BCR monitoring"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "with granular purpose specification",
          context: "Enhanced consent requirements"
        },
        {
          type: "added",
          baseText: "",
          comparisonText:
            "Data Localization Compliance: Adherence to local data residency requirements\nTransfer Impact Assessment: Mandatory evaluation of destination country privacy laws",
          context: "New transfer requirements"
        },
        {
          type: "modified",
          baseText: "must be logged",
          comparisonText: "must be logged in real-time monitoring systems",
          context: "Enhanced logging"
        },
        {
          type: "modified",
          baseText: "data mapping documentation",
          comparisonText: "comprehensive data mapping documentation with risk assessments",
          context: "Enhanced documentation"
        }
      ],
      filename: "enterprise-data-governance-policy.pdf",
      baseVersion: "1.0",
      comparisonVersion: "2.0",
      sourceTrace: "Section 3.1, Page 2",
      similarities: [
        "Maintains GDPR transfer mechanism framework",
        "Preserves adequacy decision recognition",
        "Continues requirement for transfer documentation"
      ],
      keyDifferences: [
        "Added mandatory transfer impact assessments",
        "Introduced data localization compliance requirements",
        "Enhanced SCC requirements with supplementary measures",
        "Real-time monitoring of all data transfers",
        "More granular consent and purpose specification"
      ]
    },
    {
      id: "pair-6",
      type: "aligned",
      title: "Data Security Requirements",
      section: "Section 3.2: Technical and Organizational Measures",
      baseContent:
        "All personal data must be protected through appropriate technical and organizational measures:\n\n**Technical Measures:**\n• *Encryption at rest* using AES-256\n• *Encryption in transit* using TLS 1.2+\n• *Access controls* with role-based permissions\n• *Regular security updates* and patch management\n\n**Organizational Measures:**\n• *Staff training* on data protection\n• *Incident response procedures*\n• *Regular security audits*\n• *Vendor management* processes",
      comparisonContent:
        "All personal data must be protected through robust technical and organizational measures with `continuous monitoring`:\n\n**Technical Measures:**\n• *Encryption at rest* using `AES-256 with hardware security modules`\n• *Encryption in transit* using `TLS 1.3+ with perfect forward secrecy`\n• *Zero-trust access controls* with `dynamic permissions and continuous authentication`\n• *Automated security updates* and `real-time threat detection`\n• *`Data loss prevention (DLP) systems`*\n• *`Endpoint detection and response (EDR) solutions`*\n\n**Organizational Measures:**\n• *Mandatory quarterly staff training* on data protection with `competency testing`\n• *Enhanced incident response procedures* with `24/7 response team`\n• *Continuous security monitoring* and `monthly penetration testing`\n• *Comprehensive vendor management* processes with `security scorecards`\n• *`Privacy by design implementation`* in all systems\n• *`Data protection impact assessments`* for high-risk processing",
      diffSegments: [
        {
          type: "modified",
          baseText: "appropriate technical and organizational measures",
          comparisonText: "robust technical and organizational measures with continuous monitoring",
          context: "Enhanced security approach"
        },
        {
          type: "modified",
          baseText: "AES-256",
          comparisonText: "AES-256 with hardware security modules",
          context: "Enhanced encryption"
        },
        {
          type: "modified",
          baseText: "TLS 1.2+",
          comparisonText: "TLS 1.3+ with perfect forward secrecy",
          context: "Updated transport security"
        },
        {
          type: "modified",
          baseText: "Access controls with role-based permissions",
          comparisonText: "Zero-trust access controls with dynamic permissions and continuous authentication",
          context: "Enhanced access control"
        },
        {
          type: "modified",
          baseText: "Regular security updates",
          comparisonText: "Automated security updates and real-time threat detection",
          context: "Automated security management"
        },
        {
          type: "added",
          baseText: "",
          comparisonText: "Data loss prevention (DLP) systems\nEndpoint detection and response (EDR) solutions",
          context: "New technical controls"
        },
        {
          type: "modified",
          baseText: "Staff training",
          comparisonText: "Mandatory quarterly staff training with competency testing",
          context: "Enhanced training requirements"
        },
        {
          type: "added",
          baseText: "",
          comparisonText:
            "Privacy by design implementation in all systems\nData protection impact assessments for high-risk processing",
          context: "New organizational measures"
        }
      ],
      filename: "enterprise-data-governance-policy.pdf",
      baseVersion: "1.0",
      comparisonVersion: "2.0",
      sourceTrace: "Section 3.2, Page 2",
      similarities: [
        "Maintains encryption requirements for data protection",
        "Preserves staff training and audit requirements",
        "Continues vendor management processes"
      ],
      keyDifferences: [
        "Upgraded to zero-trust security architecture",
        "Enhanced encryption with hardware security modules",
        "Added automated security monitoring and DLP systems",
        "Mandatory quarterly training with competency testing",
        "Integrated privacy by design and DPIA requirements"
      ]
    },
    {
      id: "pair-7",
      type: "unmatched_base",
      title: "Legacy Cookie Consent Framework",
      section: "Section 4.1: Deprecated Consent Mechanisms",
      baseContent:
        "Website visitors provide consent through continued browsing behavior after being presented with a cookie notice. Implied consent is considered valid for:\n\n• *Functional cookies* necessary for website operation\n• *Analytics cookies* for understanding user behavior\n• *Marketing cookies* when users don't opt-out within 30 days\n\nConsent preferences are stored in *browser local storage* and respected for 12 months before re-requesting consent.",
      filename: "enterprise-data-governance-policy.pdf",
      baseVersion: "1.0",
      sourceTrace: "Section 4.1, Page 2",
      similarities: [],
      keyDifferences: [
        "Implied consent mechanism completely removed",
        "No longer accepts continued browsing as consent",
        "Eliminates opt-out based marketing cookie consent",
        "Replaced with explicit consent requirements"
      ]
    },
    {
      id: "pair-8",
      type: "unmatched_comparison",
      title: "AI and Automated Decision-Making Governance",
      section: "Section 4.1: Algorithmic Transparency and Accountability",
      comparisonContent:
        "All AI systems and automated decision-making processes that affect individuals must comply with algorithmic governance requirements:\n\n**Transparency Requirements:**\n• *Meaningful information* about the logic, significance, and consequences of automated decisions\n• *Plain language explanations* of how AI systems make decisions\n• *Data source documentation* including training data provenance\n• *Model performance metrics* and bias testing results\n\n**Individual Rights:**\n• **Right to human review** of automated decisions within 5 business days\n• **Right to explanation** of specific decisions affecting the individual\n• **Right to contest** automated decisions with appeal process\n• **Right to opt-out** of purely automated decision-making where legally permitted\n\n**Compliance Measures:**\n• *Algorithmic impact assessments* for high-risk AI systems\n• *Continuous bias monitoring* and mitigation procedures\n• *Regular model audits* by independent third parties\n• *Incident reporting* for AI system failures or bias detection",
      filename: "enterprise-data-governance-policy.pdf",
      comparisonVersion: "2.0",
      sourceTrace: "Section 4.1, Page 2",
      similarities: [],
      keyDifferences: [
        "Comprehensive AI governance framework introduced",
        "Mandatory transparency and explainability requirements",
        "Individual rights for AI decision review and contest",
        "Continuous bias monitoring and independent audits required"
      ]
    },
    {
      id: "pair-9",
      type: "unmatched_comparison",
      title: "Incident Response and Breach Notification",
      section: "Section 4.2: Security Incident Management",
      comparisonContent:
        "Data security incidents and personal data breaches must be managed according to the following enhanced procedures:\n\n**Incident Detection and Response:**\n• *Automated detection systems* with 24/7 monitoring\n• *Immediate containment procedures* within 1 hour of detection\n• *Incident classification* based on risk severity (Critical, High, Medium, Low)\n• *Cross-functional response team* activation within 2 hours\n\n**Breach Notification Requirements:**\n• **Supervisory Authority Notification**: Within 72 hours of awareness\n• **Data Subject Notification**: Within 72 hours for high-risk breaches\n• **Internal Stakeholder Communication**: Within 4 hours including legal and executive teams\n• **Public Disclosure**: When required by law or affecting public safety\n\n**Post-Incident Activities:**\n• *Root cause analysis* within 30 days\n• *Remediation plan implementation* with timeline and accountability\n• *Lessons learned documentation* and process improvements\n• *Regulatory reporting* and compliance demonstration",
      filename: "enterprise-data-governance-policy.pdf",
      comparisonVersion: "2.0",
      sourceTrace: "Section 4.2, Page 2",
      similarities: [],
      keyDifferences: [
        "New comprehensive incident response framework",
        "Automated detection with 24/7 monitoring capability",
        "Strict timeline requirements for notification and response",
        "Enhanced post-incident analysis and improvement processes"
      ]
    },
    {
      id: "pair-10",
      type: "aligned",
      title: "Compliance Monitoring and Enforcement",
      section: "Section 5.1: Governance Oversight",
      baseContent:
        "Compliance with this policy is monitored through:\n\n• *Annual privacy audits* conducted by internal audit team\n• *Quarterly compliance reporting* to executive leadership\n• *Monthly policy training* for all staff handling personal data\n• *Vendor compliance assessments* conducted annually\n\nNon-compliance may result in disciplinary action including training, warnings, or termination for serious violations.",
      comparisonContent:
        "Compliance with this policy is continuously monitored through enhanced oversight mechanisms:\n\n• *Continuous automated compliance monitoring* with real-time dashboards\n• *Quarterly independent privacy audits* with external certification\n• *Monthly mandatory policy training* with competency testing and certification\n• *Comprehensive vendor assessments* including on-site audits and security scorecards\n• *`Data governance committee`* meetings every two weeks\n• *`Privacy officer reporting`* directly to board of directors monthly\n• *`Compliance metrics tracking`* with KPI dashboards and trend analysis\n\nNon-compliance results in progressive disciplinary action: *documented coaching*, *formal warnings*, *suspension of data access*, and *termination for serious violations*. Financial penalties may apply for `willful violations` with amounts determined by violation severity.",
      diffSegments: [
        {
          type: "modified",
          baseText: "monitored through",
          comparisonText: "continuously monitored through enhanced oversight mechanisms",
          context: "Enhanced monitoring approach"
        },
        {
          type: "modified",
          baseText: "Annual privacy audits conducted by internal audit team",
          comparisonText: "Continuous automated compliance monitoring with real-time dashboards",
          context: "Automated compliance monitoring"
        },
        {
          type: "modified",
          baseText: "Quarterly compliance reporting",
          comparisonText: "Quarterly independent privacy audits with external certification",
          context: "Independent audit approach"
        },
        {
          type: "modified",
          baseText: "Monthly policy training",
          comparisonText: "Monthly mandatory policy training with competency testing and certification",
          context: "Enhanced training requirements"
        },
        {
          type: "modified",
          baseText: "Vendor compliance assessments conducted annually",
          comparisonText: "Comprehensive vendor assessments including on-site audits and security scorecards",
          context: "Enhanced vendor oversight"
        },
        {
          type: "added",
          baseText: "",
          comparisonText:
            "Data governance committee meetings every two weeks\nPrivacy officer reporting directly to board of directors monthly\nCompliance metrics tracking with KPI dashboards and trend analysis",
          context: "New governance structures"
        },
        {
          type: "modified",
          baseText: "disciplinary action including training, warnings, or termination",
          comparisonText:
            "progressive disciplinary action: documented coaching, formal warnings, suspension of data access, and termination for serious violations. Financial penalties may apply for willful violations",
          context: "Enhanced enforcement measures"
        }
      ],
      filename: "enterprise-data-governance-policy.pdf",
      baseVersion: "1.0",
      comparisonVersion: "2.0",
      sourceTrace: "Section 5.1, Page 2",
      similarities: [
        "Maintains regular audit and reporting requirements",
        "Preserves training requirements for data handlers",
        "Continues vendor compliance assessment practices"
      ],
      keyDifferences: [
        "Introduced continuous automated monitoring",
        "Added data governance committee and board reporting",
        "Enhanced vendor assessment with on-site audits",
        "Progressive disciplinary action with financial penalties",
        "Real-time compliance dashboards and KPI tracking"
      ]
    }
  ]
}

const groupedPairs = computed(() => {
  if (!diffData.value) return { aligned: [], unmatched_base: [], unmatched_comparison: [] }

  return {
    aligned: diffData.value.pairs.filter((p) => p.type === "aligned"),
    unmatched_base: diffData.value.pairs.filter((p) => p.type === "unmatched_base"),
    unmatched_comparison: diffData.value.pairs.filter((p) => p.type === "unmatched_comparison")
  }
})

const selectedPair = computed(() => {
  if (!selectedPairId.value || !diffData.value) return null
  return diffData.value.pairs.find((p) => p.id === selectedPairId.value) || null
})

const copyRunId = async () => {
  if (diffData.value?.metadata.runId) {
    await navigator.clipboard.writeText(diffData.value.metadata.runId)
  }
}

const selectPair = (pairId: string) => {
  selectedPairId.value = pairId
}

const navigateNext = () => {
  if (!diffData.value) return
  const currentIndex = diffData.value.pairs.findIndex((p) => p.id === selectedPairId.value)
  if (currentIndex < diffData.value.pairs.length - 1) {
    selectedPairId.value = diffData.value.pairs[currentIndex + 1].id
  }
}

const navigatePrevious = () => {
  if (!diffData.value) return
  const currentIndex = diffData.value.pairs.findIndex((p) => p.id === selectedPairId.value)
  if (currentIndex > 0) {
    selectedPairId.value = diffData.value.pairs[currentIndex - 1].id
  }
}

const handleSectionSelect = (pairId: string) => {
  selectedPairId.value = pairId
  viewMode.value = "sections"
}

const setViewMode = (mode: "sections" | "overview") => {
  viewMode.value = mode
}

onMounted(async () => {
  await new Promise((resolve) => setTimeout(resolve, 1000))
  diffData.value = mockData
  selectedPairId.value = mockData.pairs[0].id
  isLoading.value = false
})

definePageMeta({
  layout: "default"
})
</script>

<template>
  <div class="h-screen flex flex-col bg-gray-50">
    <!-- <PolicyDiffHeader v-if="diffData" :metadata="diffData.metadata" @copy-run-id="copyRunId" /> -->

    <!-- View Toggle -->
    <div class="border-b border-gray-200 bg-white px-6 py-3">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-1 bg-gray-100 rounded-lg p-1">
          <button
            @click="setViewMode('overview')"
            :class="[
              'px-3 py-1.5 text-sm font-medium rounded-md transition-colors',
              viewMode === 'overview' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'
            ]"
          >
            📄 Full Document
          </button>
          <button
            @click="setViewMode('sections')"
            :class="[
              'px-3 py-1.5 text-sm font-medium rounded-md transition-colors',
              viewMode === 'sections' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'
            ]"
          >
            📝 Section by Section
          </button>
        </div>

        <div v-if="diffData" class="text-sm text-gray-600">
          {{
            groupedPairs.aligned.length + groupedPairs.unmatched_base.length + groupedPairs.unmatched_comparison.length
          }}
          sections • {{ groupedPairs.aligned.length }} modified • {{ groupedPairs.unmatched_base.length }} removed •
          {{ groupedPairs.unmatched_comparison.length }} added
        </div>
      </div>
    </div>

    <div class="flex-1 flex overflow-hidden">
      <!-- Navigation Panel (only show in sections mode) -->
      <PolicyDiffNavigator
        v-if="diffData && viewMode === 'sections'"
        :grouped-pairs="groupedPairs"
        :selected-pair-id="selectedPairId"
        @select-pair="selectPair"
        class="w-80 border-r border-gray-200"
      />

      <!-- Main Content Area -->
      <div class="flex-1 flex flex-col">
        <!-- Full Document Overview -->
        <PolicyDiffDocumentOverview
          v-if="diffData && viewMode === 'overview' && !isLoading"
          :diff-data="diffData"
          @select-section="handleSectionSelect"
          class="flex-1 h-0"
        />

        <!-- Section by Section View -->
        <PolicyDiffMarkdownViewer
          v-else-if="selectedPair && viewMode === 'sections' && !isLoading"
          :pair="selectedPair"
          @navigate-next="navigateNext"
          @navigate-previous="navigatePrevious"
          class="flex-1 h-0"
        />

        <!-- Loading State -->
        <div v-else-if="isLoading" class="flex-1 flex items-center justify-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>

        <!-- Empty State for Sections View -->
        <div v-else-if="viewMode === 'sections'" class="flex-1 flex items-center justify-center text-gray-500">
          <div class="text-center">
            <div class="text-lg mb-2">📝</div>
            <div>Select a policy section to view detailed differences</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
